package grpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/bufbuild/connect-go"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/pkg/grpc/bufcurl"
	"github.com/protoflow-labs/protoflow/pkg/grpc/protoencoding"
	"github.com/protoflow-labs/protoflow/pkg/grpc/reflect"
	"github.com/protoflow-labs/protoflow/pkg/grpc/verbose"
	"golang.org/x/net/http2"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type ReflectionManager struct {
	URL      string
	Protocol string
	Headers  []string

	resolver       protoencoding.Resolver
	requestHeaders http.Header
	httpClient     connect.HTTPClient
}

func NewReflectionManager(url string) *ReflectionManager {
	return &ReflectionManager{
		URL:      url,
		Protocol: connect.ProtocolGRPC,
	}
}

type ReflectionOptions struct {
	Protocol string
	Headers  []string
}

type InvokeOptions struct {
	OutputStream    bufcurl.OutputStream
	TLSConfig       bufcurl.TLSSettings
	URL             string
	Protocol        string
	Headers         []string
	UserAgent       string
	ReflectProtocol string
	ReflectHeaders  []string

	// Protocol details
	UnixSocket          string
	HTTP2PriorKnowledge bool

	// Timeouts
	NoKeepAlive           bool
	KeepAliveTimeSeconds  float64
	ConnectTimeoutSeconds float64
}

func makeHTTPClient(f InvokeOptions, isSecure bool, authority string, printer verbose.Printer) (connect.HTTPClient, error) {
	var dialer net.Dialer
	if f.ConnectTimeoutSeconds != 0 {
		dialer.Timeout = secondsToDuration(f.ConnectTimeoutSeconds)
	}
	if f.NoKeepAlive {
		dialer.KeepAlive = -1
	} else {
		dialer.KeepAlive = secondsToDuration(f.KeepAliveTimeSeconds)
	}
	var dialFunc func(ctx context.Context, network, address string) (net.Conn, error)
	if f.UnixSocket != "" {
		dialFunc = func(ctx context.Context, network, addr string) (net.Conn, error) {
			printer.Printf("* Dialing unix socket %s...", f.UnixSocket)
			return dialer.DialContext(ctx, "unix", f.UnixSocket)
		}
	} else {
		dialFunc = func(ctx context.Context, network, address string) (net.Conn, error) {
			printer.Printf("* Dialing (%s) %s...", network, address)
			conn, err := dialer.DialContext(ctx, network, address)
			if err != nil {
				return nil, err
			}
			printer.Printf("* Connected to %s", conn.RemoteAddr().String())
			return conn, err
		}
	}

	var transport http.RoundTripper
	if !isSecure && f.HTTP2PriorKnowledge {
		transport = &http2.Transport{
			AllowHTTP: true,
			DialTLSContext: func(ctx context.Context, network, addr string, _ *tls.Config) (net.Conn, error) {
				return dialFunc(ctx, network, addr)
			},
		}
	} else {
		var tlsConfig *tls.Config
		if isSecure {
			var err error
			tlsConfig, err = bufcurl.MakeVerboseTLSConfig(&f.TLSConfig, authority, printer)
			if err != nil {
				return nil, err
			}
		}
		transport = &http.Transport{
			Proxy:             http.ProxyFromEnvironment,
			DialContext:       dialFunc,
			ForceAttemptHTTP2: true,
			MaxIdleConns:      1,
			TLSClientConfig:   tlsConfig,
		}
	}
	return bufcurl.NewVerboseHTTPClient(transport, printer), nil
}

func Test(f InvokeOptions, input bufcurl.InputStream) error {
	// TODO breadchris replace container with zerolog, we don't need to be passing this around rn
	container := &bufcurl.Container{}

	endpointURL, service, method, baseURL, err := verifyEndpointURL(f.URL)
	if err != nil {
		return err
	}
	isSecure := endpointURL.Scheme == "https"

	var clientOptions []connect.ClientOption
	switch f.Protocol {
	case connect.ProtocolGRPC:
		clientOptions = []connect.ClientOption{connect.WithGRPC()}
	case connect.ProtocolGRPCWeb:
		clientOptions = []connect.ClientOption{connect.WithGRPCWeb()}
	}
	if f.Protocol != connect.ProtocolGRPC {
		// The transport will log trailers to the verbose printer. But if
		// we're not using standard grpc protocol, trailers are actually encoded
		// in an end-of-stream message for streaming calls. So this interceptor
		// will print the trailers for streaming calls when the response stream
		// is drained.
		clientOptions = append(clientOptions, connect.WithInterceptors(bufcurl.TraceTrailersInterceptor(container.VerbosePrinter())))
	}

	requestHeaders, err := bufcurl.LoadHeaders(f.Headers)
	if err != nil {
		return err
	}
	if len(requestHeaders.Values("user-agent")) == 0 {
		userAgent := f.UserAgent
		if userAgent == "" {
			userAgent = bufcurl.DefaultUserAgent(f.Protocol, "1.0.0")
		}
		requestHeaders.Set("user-agent", userAgent)
	}

	transport, err := makeHTTPClient(f, isSecure, bufcurl.GetAuthority(endpointURL, requestHeaders), container.VerbosePrinter())
	if err != nil {
		return err
	}

	reflectHeaders, err := bufcurl.LoadHeaders(f.ReflectHeaders)
	if err != nil {
		return err
	}
	reflectProtocol, err := bufcurl.ParseReflectProtocol(f.ReflectProtocol)
	if err != nil {
		return err
	}
	var closeRes func()
	res, closeRes := bufcurl.NewServerReflectionResolver(context.Background(), transport, clientOptions, baseURL, reflectProtocol, reflectHeaders, container.VerbosePrinter())
	defer closeRes()
	methodDescriptor, err := bufcurl.ResolveMethodDescriptor(res, service, method)
	if err != nil {
		return err
	}
	println(methodDescriptor)

	// Now we can finally issue the RPC
	if f.OutputStream == nil {
		invoker := bufcurl.NewInvoker(methodDescriptor, res, transport, clientOptions, f.URL, os.Stdout)
		return invoker.Invoke(context.Background(), input, requestHeaders)
	} else {
		invoker := bufcurl.NewInvoker(methodDescriptor, res, transport, clientOptions, f.URL, os.Stdout)
		return invoker.InvokeWithStream(context.Background(), input, f.OutputStream, requestHeaders)
	}
}

func (s *ReflectionManager) Init() (func(), error) {
	reflectOpts := &ReflectionOptions{
		Protocol: "grpc-v1",
	}
	// TODO breadchris options

	printer := &bufcurl.ZeroLogPrinter{}

	endpointURL, err := verifyServerEndpointURL(s.URL)
	if err != nil {
		return nil, err
	}
	isSecure := endpointURL.Scheme == "https"

	s.requestHeaders, err = loadRequestHeaders(s.Protocol, s.Headers)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load request headers")
	}

	client := reflect.NewClientBuilder()
	client.IsSecure = isSecure
	client.Authority = bufcurl.GetAuthority(endpointURL, s.requestHeaders)

	s.httpClient, err = client.Build()
	if err != nil {
		return nil, err
	}

	reflectHeaders, err := bufcurl.LoadHeaders(reflectOpts.Headers)
	if err != nil {
		return nil, err
	}
	reflectProtocol, err := bufcurl.ParseReflectProtocol(reflectOpts.Protocol)
	if err != nil {
		return nil, err
	}

	clientOptions := loadClientOptions(s.Protocol)
	var cleanup func()
	s.resolver, cleanup = bufcurl.NewServerReflectionResolver(
		context.Background(),
		s.httpClient,
		clientOptions,
		endpointURL.String(),
		reflectProtocol,
		reflectHeaders,
		printer,
	)
	return cleanup, nil
}

func (s *ReflectionManager) ExecuteMethod(
	ctx context.Context,
	method protoreflect.MethodDescriptor,
	input bufcurl.InputStream,
	output bufcurl.OutputStream,
) (err error) {
	grpcURL := fmt.Sprintf("%s/%s/%s", s.URL, method.Parent().FullName(), method.Name())
	clientOptions := loadClientOptions(s.Protocol)

	invoker := bufcurl.NewInvoker(method, s.resolver, s.httpClient, clientOptions, grpcURL, os.Stdout)
	return invoker.InvokeWithStream(ctx, input, output, s.requestHeaders)
}

func (s *ReflectionManager) ResolveMethod(service, method string) (protoreflect.MethodDescriptor, error) {
	descriptor, err := s.resolver.FindDescriptorByName(protoreflect.FullName(service))
	if err == protoregistry.NotFound {
		return nil, errors.Wrapf(err, "failed to find service in schema")
	} else if err != nil {
		return nil, errors.Wrapf(err, "connection error when resolving method")
	}
	serviceDescriptor, ok := descriptor.(protoreflect.ServiceDescriptor)
	if !ok {
		return nil, errors.Wrapf(err, "failed to find method")
	}
	methodDescriptor := serviceDescriptor.Methods().ByName(protoreflect.Name(method))
	if methodDescriptor == nil {
		return nil, errors.Wrapf(err, "failed to find method in service descriptor")
	}
	return methodDescriptor, nil
}

func verifyServerEndpointURL(urlArg string) (endpointURL *url.URL, err error) {
	endpointURL, err = url.Parse(urlArg)
	if err != nil {
		return nil, fmt.Errorf("%q is not a valid endpoint URL: %w", urlArg, err)
	}
	if endpointURL.Scheme != "http" && endpointURL.Scheme != "https" {
		return nil, fmt.Errorf("invalid endpoint URL: scheme %q is not supported", endpointURL.Scheme)
	}

	if strings.HasSuffix(endpointURL.Path, "/") {
		return nil, fmt.Errorf("invalid endpoint URL: path %q should not end with a slash (/)", endpointURL.Path)
	}
	return endpointURL, nil
}

func loadRequestHeaders(protocol string, headers []string) (http.Header, error) {
	requestHeaders, err := bufcurl.LoadHeaders(headers)
	if err != nil {
		return nil, err
	}
	if len(requestHeaders.Values("user-agent")) == 0 {
		userAgent := bufcurl.DefaultUserAgent(protocol, "1.0.0")
		requestHeaders.Set("user-agent", userAgent)
	}
	return requestHeaders, nil
}

func loadClientOptions(protocol string) []connect.ClientOption {
	var clientOptions []connect.ClientOption

	switch protocol {
	case connect.ProtocolGRPCWeb:
		clientOptions = []connect.ClientOption{connect.WithGRPCWeb()}
	case connect.ProtocolGRPC:
		clientOptions = []connect.ClientOption{connect.WithGRPC()}
	}

	if protocol != connect.ProtocolGRPC {
		// The transport will log trailers to the verbose printer. But if
		// we're not using standard grpc protocol, trailers are actually encoded
		// in an end-of-stream message for streaming calls. So this interceptor
		// will print the trailers for streaming calls when the response stream
		// is drained.
		clientOptions = append(clientOptions, connect.WithInterceptors(
			bufcurl.TraceTrailersInterceptor(&bufcurl.ZeroLogPrinter{}),
		))
	}
	return clientOptions
}
