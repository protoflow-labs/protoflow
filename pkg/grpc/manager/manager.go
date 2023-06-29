package manager

import (
	"context"
	"fmt"
	"github.com/bufbuild/connect-go"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/pkg/grpc/bufcurl"
	"github.com/protoflow-labs/protoflow/pkg/grpc/bufcurl/protoencoding"
	"github.com/protoflow-labs/protoflow/pkg/grpc/bufcurl/reflect"
	"github.com/protoflow-labs/protoflow/pkg/util/rx"
	"github.com/reactivex/rxgo/v2"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
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
	reflectOps     ReflectionOptions
}

type ReflectionOptions struct {
	Protocol bufcurl.ReflectProtocol
	Headers  []string
}

type Options func(*ReflectionManager)

func WithProtocol(protocol bufcurl.ReflectProtocol) Options {
	return func(m *ReflectionManager) {
		m.reflectOps.Protocol = protocol
	}
}

func NewReflectionManager(url string, ops ...Options) *ReflectionManager {
	m := &ReflectionManager{
		URL:      url,
		Protocol: connect.ProtocolGRPC,
	}
	for _, op := range ops {
		op(m)
	}
	return m
}

func (s *ReflectionManager) Init() (func(), error) {
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
	reflectHeaders, err := bufcurl.LoadHeaders(s.reflectOps.Headers)
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
		s.reflectOps.Protocol,
		reflectHeaders,
		printer,
	)
	return cleanup, nil
}

func (s *ReflectionManager) ExecuteMethod(
	ctx context.Context,
	method protoreflect.MethodDescriptor,
	input rxgo.Observable,
	output rx.ItemSink,
) (err error) {
	grpcURL := fmt.Sprintf("%s/%s/%s", s.URL, method.Parent().FullName(), method.Name())
	clientOptions := loadClientOptions(s.Protocol)

	invoker := bufcurl.NewInvoker(method, s.resolver, s.httpClient, clientOptions, grpcURL, os.Stdout, output)
	return invoker.Invoke(ctx, input, s.requestHeaders)
}

func (s *ReflectionManager) ResolveMethod(service, methodName string) (protoreflect.MethodDescriptor, error) {
	descriptor, err := s.resolver.FindDescriptorByName(protoreflect.FullName(service))
	if err == protoregistry.NotFound {
		return nil, errors.Wrapf(err, "failed to find service in schema")
	} else if err != nil {
		return nil, errors.Wrapf(err, "connection error when resolving methodName")
	}
	serviceDescriptor, ok := descriptor.(protoreflect.ServiceDescriptor)
	if !ok {
		return nil, errors.Wrapf(err, "failed to find methodName")
	}

	methodDescriptor := serviceDescriptor.Methods().ByName(protoreflect.Name(methodName))
	if methodDescriptor == nil {

		return nil, errors.New("failed to find methodName in service descriptor")
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
