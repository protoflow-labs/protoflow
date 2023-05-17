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
	"github.com/protoflow-labs/protoflow/pkg/grpc/verbose"
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

func ExecuteCurl(ctx context.Context, f InvokeOptions, input bufcurl.InputStream) (err error) {
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
	res, closeRes := bufcurl.NewServerReflectionResolver(ctx, transport, clientOptions, baseURL, reflectProtocol, reflectHeaders, container.VerbosePrinter())
	defer closeRes()

	methodDescriptor, err := bufcurl.ResolveMethodDescriptor(res, service, method)
	if err != nil {
		return err
	}

	// Now we can finally issue the RPC
	if f.OutputStream == nil {
		invoker := bufcurl.NewInvoker(container, methodDescriptor, res, transport, clientOptions, f.URL, os.Stdout)
		return invoker.Invoke(ctx, input, requestHeaders)
	} else {
		invoker := bufcurl.NewInvoker(container, methodDescriptor, res, transport, clientOptions, f.URL, os.Stdout)
		return invoker.InvokeWithStream(ctx, input, f.OutputStream, requestHeaders)
	}
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

func secondsToDuration(secs float64) time.Duration {
	return time.Duration(float64(time.Second) * secs)
}
