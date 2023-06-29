package resource

import (
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/workflow/node"
	"github.com/rs/zerolog/log"
	"net"
	"net/url"
	"strings"
	"time"
)

type LanguageServiceResource struct {
	*gen.LanguageService
	*GRPCResource
}

func (r *LanguageServiceResource) Name() string {
	return LanguageServiceType
}

func (r *LanguageServiceResource) Init() (func(), error) {
	return r.GRPCResource.Init()
}

func (r *LanguageServiceResource) Info(n node.Node) (*node.Info, error) {
	f, ok := n.(*node.FunctionNode)
	if !ok {
		return nil, errors.New("node is not a function node")
	}
	grpcNode := r.ToGRPC(f)
	// TODO breadchris we should know where the function node is located and should read/write from the proto
	return r.GRPCResource.Info(grpcNode)
}

func (r *LanguageServiceResource) ToGRPC(n *node.FunctionNode) *node.GRPCNode {
	serviceName := strings.ToLower(r.Runtime.String()) + "Service"
	return &node.GRPCNode{
		BaseNode: n.BaseNode,
		GRPC: &gen.GRPC{
			Package: "protoflow",
			Service: serviceName,
			Method:  n.Name,
		},
	}
}

func ensureRunning(host string) error {
	maxRetries := 1
	retryInterval := 2 * time.Second

	u, err := url.Parse(host)
	if err != nil {
		return errors.Wrapf(err, "unable to parse url %s", host)
	}

	log.Debug().Str("host", host).Msg("waiting for host to come online")
	for i := 1; i <= maxRetries; i++ {
		conn, err := net.DialTimeout("tcp", u.Host, time.Second)
		if err == nil {
			conn.Close()
			log.Debug().Str("host", host).Msg("host is not listening")
			return nil
		} else {
			log.Debug().Err(err).Int("attempt", i).Int("max", maxRetries).Msg("error connecting to host")
			time.Sleep(retryInterval)
		}
	}
	return errors.New("host did not come online in time")
}
