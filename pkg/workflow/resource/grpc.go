package resource

import (
	"github.com/protoflow-labs/protoflow/gen"
	"strings"
)

type GRPCResource struct {
	*BaseResource
	*gen.GRPCService
}

func (r *GRPCResource) Init() (func(), error) {
	// TODO breadchris this is a hack to get the grpc server running, this is not ideal
	if !strings.HasPrefix(r.Host, "http://") {
		r.Host = "http://" + r.Host
	}
	//if err := ensureRunning(r.Host); err != nil {
	//	// TODO breadchris ignore errors for now
	//	// return nil, errors.Wrapf(err, "unable to get the %s grpc server running", r.Name())
	//	return nil, nil
	//}
	return nil, nil
}
