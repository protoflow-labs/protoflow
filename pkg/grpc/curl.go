// Copyright 2020-2023 Buf Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package grpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/protoflow-labs/protoflow/pkg/grpc/bufcurl"
	"github.com/protoflow-labs/protoflow/pkg/grpc/protoencoding"
	"github.com/protoflow-labs/protoflow/pkg/grpc/verbose"
	"google.golang.org/protobuf/reflect/protoreflect"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/bufbuild/connect-go"
	"golang.org/x/net/http2"
)

func verifyEndpointURL(urlArg string) (endpointURL *url.URL, service, method, baseURL string, err error) {
	endpointURL, err = url.Parse(urlArg)
	if err != nil {
		return nil, "", "", "", fmt.Errorf("%q is not a valid endpoint URL: %w", urlArg, err)
	}
	if endpointURL.Scheme != "http" && endpointURL.Scheme != "https" {
		return nil, "", "", "", fmt.Errorf("invalid endpoint URL: scheme %q is not supported", endpointURL.Scheme)
	}

	if strings.HasSuffix(endpointURL.Path, "/") {
		return nil, "", "", "", fmt.Errorf("invalid endpoint URL: path %q should not end with a slash (/)", endpointURL.Path)
	}
	parts := strings.Split(endpointURL.Path, "/")
	if len(parts) < 2 || parts[len(parts)-1] == "" || parts[len(parts)-2] == "" {
		return nil, "", "", "", fmt.Errorf("invalid endpoint URL: path %q should end with two non-empty components indicating service and method", endpointURL.Path)
	}
	service, method = parts[len(parts)-2], parts[len(parts)-1]
	baseURL = strings.TrimSuffix(urlArg, service+"/"+method)
	if baseURL == urlArg {
		// should not be possible due to above checks
		return nil, "", "", "", fmt.Errorf("failed to extract base URL from %q", urlArg)
	}
	return endpointURL, service, method, baseURL, nil
}

type RemoteMethod struct {
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

	// TODO breadchris values are stored here so we can get the method descriptor
	MethodDescriptor protoreflect.MethodDescriptor
	resolver         protoencoding.Resolver
	closeResolver    func()
	transport        connect.HTTPClient
	clientOptions    []connect.ClientOption
	requestHeaders   http.Header
}

// TODO breadchris we should passing options to this function instead of having a bunch of fields
func NewRemoteMethod(g *RemoteMethod) (*RemoteMethod, error) {
	err := g.init()
	if err != nil {
		return nil, err
	}
	return g, nil
}

// TODO breadchris spend more than 2 brain cells organizing this code
func (g *RemoteMethod) init() error {
	printer := &bufcurl.ZeroLogPrinter{}

	endpointURL, service, method, baseURL, err := verifyEndpointURL(g.URL)
	if err != nil {
		return err
	}
	isSecure := endpointURL.Scheme == "https"

	switch g.Protocol {
	case connect.ProtocolGRPC:
		g.clientOptions = []connect.ClientOption{connect.WithGRPC()}
	case connect.ProtocolGRPCWeb:
		g.clientOptions = []connect.ClientOption{connect.WithGRPCWeb()}
	}
	if g.Protocol != connect.ProtocolGRPC {
		// The transport will log trailers to the verbose printer. But if
		// we're not using standard grpc protocol, trailers are actually encoded
		// in an end-of-stream message for streaming calls. So this interceptor
		// will print the trailers for streaming calls when the response stream
		// is drained.
		g.clientOptions = append(g.clientOptions, connect.WithInterceptors(bufcurl.TraceTrailersInterceptor(printer)))
	}

	g.requestHeaders, err = bufcurl.LoadHeaders(g.Headers)
	if err != nil {
		return err
	}
	if len(g.requestHeaders.Values("user-agent")) == 0 {
		userAgent := g.UserAgent
		if userAgent == "" {
			userAgent = bufcurl.DefaultUserAgent(g.Protocol, "1.0.0")
		}
		g.requestHeaders.Set("user-agent", userAgent)
	}

	g.transport, err = g.makeHTTPClient(isSecure, bufcurl.GetAuthority(endpointURL, g.requestHeaders), printer)
	if err != nil {
		return err
	}

	reflectHeaders, err := bufcurl.LoadHeaders(g.ReflectHeaders)
	if err != nil {
		return err
	}
	reflectProtocol, err := bufcurl.ParseReflectProtocol(g.ReflectProtocol)
	if err != nil {
		return err
	}
	g.resolver, g.closeResolver = bufcurl.NewServerReflectionResolver(
		context.Background(), g.transport, g.clientOptions, baseURL, reflectProtocol, reflectHeaders, printer)

	g.MethodDescriptor, err = bufcurl.ResolveMethodDescriptor(g.resolver, service, method)
	if err != nil {
		return err
	}
	return nil
}

func (g *RemoteMethod) Execute(ctx context.Context, input bufcurl.InputStream) (err error) {
	defer g.closeResolver()

	// TODO breadchris what is the difference between a method and an invoker?
	invoker := bufcurl.NewInvoker(g.MethodDescriptor, g.resolver, g.transport, g.clientOptions, g.URL, os.Stdout)
	return invoker.InvokeWithStream(ctx, input, g.OutputStream, g.requestHeaders)
}

func (g *RemoteMethod) makeHTTPClient(isSecure bool, authority string, printer verbose.Printer) (connect.HTTPClient, error) {
	var dialer net.Dialer
	if g.ConnectTimeoutSeconds != 0 {
		dialer.Timeout = secondsToDuration(g.ConnectTimeoutSeconds)
	}
	if g.NoKeepAlive {
		dialer.KeepAlive = -1
	} else {
		dialer.KeepAlive = secondsToDuration(g.KeepAliveTimeSeconds)
	}
	var dialFunc func(ctx context.Context, network, address string) (net.Conn, error)
	if g.UnixSocket != "" {
		dialFunc = func(ctx context.Context, network, addr string) (net.Conn, error) {
			printer.Printf("* Dialing unix socket %s...", g.UnixSocket)
			return dialer.DialContext(ctx, "unix", g.UnixSocket)
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
	if !isSecure && g.HTTP2PriorKnowledge {
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
			tlsConfig, err = bufcurl.MakeVerboseTLSConfig(&g.TLSConfig, authority, printer)
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

func secondsToDuration(secs float64) time.Duration {
	return time.Duration(float64(time.Second) * secs)
}
