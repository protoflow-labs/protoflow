package workflow

import (
	"github.com/jhump/protoreflect/desc"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/grpc"
	"github.com/protoflow-labs/protoflow/pkg/grpc/manager"
)

type NodeInfo struct {
	MethodProto string
	TypeInfo    *gen.GRPCTypeInfo
}

func GetNodeInfo(project *gen.Project, nodeId string) (*NodeInfo, error) {
	node := findNode(project, nodeId)
	if node == nil {
		return nil, errors.Errorf("node not found")
	}

	w, err := FromProject(project)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get workflow from project")
	}

	nodeResources, err := getNodeResources(node, project.Resources)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get node resources")
	}

	var resp *NodeInfo
	// TODO breadchris we are probably going to want to reuse this logic, so there should probably be a NodeBuilder?
	switch node.Config.(type) {
	case *gen.Node_Input:
		children := w.AdjMap[nodeId]
		if len(children) != 1 {
			// TODO breadchris support multiple children
			return nil, errors.Errorf("input node should have 1 child, got %d", len(children))
		}
		// TODO breadchris optimized for specific case
		for child := range children {
			n := findNode(project, child)
			if n == nil {
				return nil, errors.Errorf("node %s not found", child)
			}

			r, err := getNodeResources(n, project.Resources)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to get node resources")
			}
			resp, err = getNodeGRPCInfo(n.GetGrpc(), r)
			if err != nil {
				return nil, err
			}
		}
	case *gen.Node_Function:
		f := node.GetFunction()
		resp, err = getNodeGRPCInfo(f.Grpc, nodeResources)
		if err != nil {
			return nil, err
		}
	case *gen.Node_Grpc:
		resp, err = getNodeGRPCInfo(node.GetGrpc(), nodeResources)
		if err != nil {
			return nil, err
		}
	}
	return resp, nil
}

// TODO breadchris should this functionality exist on a block? A block would produce its own type information
func getNodeGRPCInfo(node *gen.GRPC, resources []*gen.Resource) (*NodeInfo, error) {
	var grpcResource *gen.GRPCService
	for _, r := range resources {
		switch r.Type.(type) {
		case *gen.Resource_LanguageService:
			grpcResource = r.GetLanguageService().Grpc
		case *gen.Resource_GrpcService:
			grpcResource = r.GetGrpcService()
		}
	}

	if grpcResource == nil {
		return nil, errors.Errorf("node does not have a grpc resource")
	}

	// TODO breadchris I think a grpc resource should have a host that has a protocol
	m := manager.NewReflectionManager("http://" + grpcResource.Host)
	cleanup, err := m.Init()
	if err != nil {
		return nil, errors.Wrapf(err, "error initializing reflection manager")
	}
	defer cleanup()

	serviceName := node.Package + "." + node.Service
	method, err := m.ResolveMethod(serviceName, node.Method)
	if err != nil {
		return nil, errors.Wrapf(err, "error resolving method")
	}

	methodProto, err := manager.GetProtoForMethod(node.Package, node.Service, method)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting proto for method")
	}

	descMethod, err := desc.WrapMethod(method)
	md := grpc.NewMethodDescriptor(descMethod.GetInputType())
	typeInfo := &gen.GRPCTypeInfo{
		Input:      descMethod.GetInputType().AsDescriptorProto(),
		Output:     descMethod.GetOutputType().AsDescriptorProto(),
		DescLookup: md.DescLookup,
		EnumLookup: md.EnumLookup,
		MethodDesc: descMethod.AsMethodDescriptorProto(),
	}

	return &NodeInfo{
		MethodProto: methodProto,
		TypeInfo:    typeInfo,
	}, nil
}

func getNodeResources(node *gen.Node, resources []*gen.Resource) ([]*gen.Resource, error) {
	var nodeResources []*gen.Resource
	for _, r := range resources {
		for _, id := range node.ResourceIds {
			if r.Id == id {
				nodeResources = append(nodeResources, r)
			}
		}
	}
	return nodeResources, nil
}

func findNode(project *gen.Project, id string) *gen.Node {
	for _, n := range project.Graph.Nodes {
		if n.Id == id {
			return n
		}
	}
	return nil
}
