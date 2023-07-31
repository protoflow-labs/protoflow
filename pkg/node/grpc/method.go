package grpc

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	pgrpc "github.com/protoflow-labs/protoflow/gen/grpc"
	"github.com/protoflow-labs/protoflow/pkg/grpc"
	"github.com/protoflow-labs/protoflow/pkg/grpc/bufcurl"
	"github.com/protoflow-labs/protoflow/pkg/grpc/manager"
	"github.com/protoflow-labs/protoflow/pkg/node/base"
	"github.com/protoflow-labs/protoflow/pkg/util/rx"
	"github.com/protoflow-labs/protoflow/pkg/workflow/graph"
	"github.com/reactivex/rxgo/v2"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/reflect/protoreflect"
	"net/url"
	"strings"
)

type Method struct {
	*base.Node
	*pgrpc.Method
}

var _ graph.Node = &Method{}

func NewMethod(b *base.Node, n *pgrpc.Method) *Method {
	return &Method{
		Node:   b,
		Method: n,
	}
}

func NewMethodProto(packageService, m string) *pgrpc.GRPC {
	ps := strings.Split(packageService, ".")
	if len(ps) != 2 {
		return nil
	}
	return &pgrpc.GRPC{
		Type: &pgrpc.GRPC_Method{
			Method: &pgrpc.Method{
				Package: ps[0],
				Service: ps[1],
				Method:  m,
			},
		},
	}
}

func GetMethodFromServer(r *Server, n *Method, protocol bufcurl.ReflectProtocol) (protoreflect.MethodDescriptor, error) {
	// TODO breadchris I think a grpc resource should have a host that has a protocol
	m := manager.NewReflectionManager("http://"+r.Host, manager.WithProtocol(protocol))
	cleanup, err := m.Init()
	if err != nil {
		return nil, errors.Wrapf(err, "error initializing reflection manager")
	}
	defer cleanup()

	serviceName := n.Package + "." + n.Service
	method, err := m.ResolveMethod(serviceName, n.Method.Method)
	if err != nil {
		return nil, errors.Wrapf(err, "error resolving method")
	}
	return method, nil
}

func GetMethodDescriptor(s *Server, n *Method) (*grpc.MethodDescriptor, error) {
	method, err := GetMethodFromServer(s, n, bufcurl.ReflectProtocolGRPCV1Alpha)
	if err != nil {
		// TODO breadchris is there a cleaner way to determine if the server supports v1?
		method, err = GetMethodFromServer(s, n, bufcurl.ReflectProtocolGRPCV1)
		if err != nil {
			return nil, errors.Wrapf(err, "error getting method from server")
		}
	}

	md, err := grpc.NewMethodDescriptor(method)
	if err != nil {
		return nil, errors.Wrapf(err, "error creating method descriptor")
	}
	return md, nil
}

func (n *Method) Info() (*graph.Info, error) {
	// TODO breadchris what if we want to get the proto from a proto file?

	p, err := n.Provider()
	if err != nil {
		return nil, errors.Wrapf(err, "error getting provider")
	}
	gr, ok := p.(*Server)
	if !ok {
		return nil, fmt.Errorf("provider is not a grpc server")
	}
	md, err := GetMethodDescriptor(gr, n)
	return &graph.Info{
		Method: md,
	}, nil
}

func formatHost(host string) (string, error) {
	u, err := url.Parse(host)
	if err != nil {
		return "", errors.Wrapf(err, "error parsing url: %s", host)
	}
	if u.Scheme == "" {
		u.Scheme = "http"
	}
	return u.String(), nil
}

// TODO breadchris this should be workflow.Context, but for the memory executor it needs context.Context
func (n *Method) Wire(ctx context.Context, input graph.Input) (graph.Output, error) {
	log.Info().
		Str("service", n.Service).
		Str("method", n.Method.Method).
		Msg("setting up grpc node")

	p, err := n.Provider()
	if err != nil {
		return graph.Output{}, errors.Wrapf(err, "error getting provider")
	}

	g, ok := p.(*Server)
	if !ok {
		return graph.Output{}, fmt.Errorf("error getting GRPC resource: %s.%s", n.Service, n.Method)
	}

	serviceName := n.Service
	if n.Package != "" {
		serviceName = n.Package + "." + serviceName
	}

	host, err := formatHost(g.Host)
	if err != nil {
		return graph.Output{}, errors.Wrapf(err, "error formatting host: %s", g.Host)
	}

	manager := manager.NewReflectionManager(host)

	cleanup, err := manager.Init()
	if err != nil {
		return graph.Output{}, errors.Wrapf(err, "error initializing reflection manager")
	}
	defer cleanup()

	method, err := manager.ResolveMethod(serviceName, n.Method.Method)
	if err != nil {
		return graph.Output{}, errors.Wrapf(err, "error resolving method: %s.%s", serviceName, n.Method)
	}

	outputStream := make(chan rxgo.Item)
	// TODO breadchris we are relying on this grpc call to close the output stream. How can the stream be closed by the caller?
	if !method.IsStreamingClient() {
		// if the method is not a client stream, we need to send each input observable as a single request
		// TODO breadchris type of this method should be inferred when the workflow is parsed
		input.Observable.ForEach(func(item any) {
			log.Debug().
				Str("name", n.NormalizedName()).
				Interface("item", item).
				Msg("executing single grpc method")

			err = manager.ExecuteMethod(ctx, method, rx.FromValues(item), outputStream)
			if err != nil {
				outputStream <- rx.NewError(errors.Wrapf(err, "error calling grpc method: %s", host))
			}
			log.Debug().
				Str("name", n.NormalizedName()).
				Interface("item", item).
				Msg("done executing single grpc method")
		}, func(err error) {
			outputStream <- rx.NewError(err)
		}, func() {
			close(outputStream)
		})
	} else {
		go func() {
			log.Debug().
				Str("name", n.NormalizedName()).
				Msg("executing streaming grpc method")
			defer close(outputStream)
			err = manager.ExecuteMethod(ctx, method, input.Observable, outputStream)
			if err != nil {
				outputStream <- rx.NewError(errors.Wrapf(err, "error calling grpc method: %s", host))
			}
		}()
	}
	return graph.Output{
		Observable: rxgo.FromChannel(outputStream, rxgo.WithPublishStrategy()),
	}, nil
}
