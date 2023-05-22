package reflect

import (
	"context"
	"crypto/tls"
	"github.com/bufbuild/connect-go"
	"github.com/protoflow-labs/protoflow/pkg/grpc/bufcurl"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/http2"
	"net"
	"net/http"
	"time"
)

type ClientBuilder struct {
	IsSecure              bool
	Authority             string
	ConnectTimeoutSeconds float64
	NoKeepAlive           bool
	KeepAliveTimeSeconds  float64
	UnixSocket            string
	HTTP2PriorKnowledge   bool
	TLSConfig             bufcurl.TLSSettings
}

func NewClientBuilder() *ClientBuilder {
	return &ClientBuilder{
		IsSecure:              false,
		Authority:             "",
		ConnectTimeoutSeconds: 0,
		NoKeepAlive:           false,
		KeepAliveTimeSeconds:  0,
		UnixSocket:            "",
		HTTP2PriorKnowledge:   true,
	}
}

func (r *ClientBuilder) Build() (connect.HTTPClient, error) {
	// TODO breadchris remove
	printer := &bufcurl.ZeroLogPrinter{}

	var dialer net.Dialer
	if r.ConnectTimeoutSeconds != 0 {
		dialer.Timeout = secondsToDuration(r.ConnectTimeoutSeconds)
	}
	if r.NoKeepAlive {
		dialer.KeepAlive = -1
	} else {
		dialer.KeepAlive = secondsToDuration(r.KeepAliveTimeSeconds)
	}
	var dialFunc func(ctx context.Context, network, address string) (net.Conn, error)
	if r.UnixSocket != "" {
		dialFunc = func(ctx context.Context, network, addr string) (net.Conn, error) {
			log.Debug().Msgf("* Dialing unix socket %s...", r.UnixSocket)
			return dialer.DialContext(ctx, "unix", r.UnixSocket)
		}
	} else {
		dialFunc = func(ctx context.Context, network, address string) (net.Conn, error) {
			log.Debug().Msgf("* Dialing (%s) %s...", network, address)
			conn, err := dialer.DialContext(ctx, network, address)
			if err != nil {
				return nil, err
			}
			log.Debug().Msgf("* Connected to %s", conn.RemoteAddr().String())
			return conn, err
		}
	}

	var transport http.RoundTripper
	if !r.IsSecure && r.HTTP2PriorKnowledge {
		transport = &http2.Transport{
			AllowHTTP: true,
			DialTLSContext: func(ctx context.Context, network, addr string, _ *tls.Config) (net.Conn, error) {
				return dialFunc(ctx, network, addr)
			},
		}
	} else {
		var tlsConfig *tls.Config
		if r.IsSecure {
			var err error
			tlsConfig, err = bufcurl.MakeVerboseTLSConfig(&r.TLSConfig, r.Authority, printer)
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
