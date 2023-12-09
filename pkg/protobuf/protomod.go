package protobuf

import (
	"fmt"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/grpc"
	"github.com/protoflow-labs/protoflow/pkg/grpc/manager"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
	"os"
	"os/exec"
	"path"
	"strings"
)

func removeLineContainingString(input, substring string) string {
	lines := strings.Split(input, "\n")
	var result []string

	for _, line := range lines {
		if !strings.Contains(line, substring) {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}

func writeProtoFile(filepath string, fd *desc.FileDescriptor) error {
	_, protoName := path.Split(filepath)

	file, err := manager.PrintFile(fd)
	if err != nil {
		return errors.Wrapf(err, "error printing file")
	}

	b, err := os.OpenFile(filepath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return errors.Wrapf(err, "error opening file")
	}

	// TODO breadchris explore why there is an import for the protofile itself
	file = removeLineContainingString(file, protoName)

	_, err = b.WriteString(file)
	b.Close()
	if err != nil {
		return errors.Wrapf(err, "error writing file")
	}
	return nil
}

func AddMethod(dir, protofile, protoPackage, service, method string) error {
	fileDescs, err := grpc.ParseProtoDir(dir, protofile)
	if err != nil {
		return errors.Wrapf(err, "error parsing protos")
	}

	// TODO breadchris when does this fail? what if the strings have problematic characters?
	fullname := protoPackage + "." + service

	var handlers []string
	for _, fd := range fileDescs {
		svc := fd.FindService(fullname)
		if svc == nil {
			log.Warn().Msgf("service not found in proto %s", fd.GetName())
			continue
		}

		// TODO breadchris how can you parse this with unmarshal?
		opts := svc.AsServiceDescriptorProto().GetOptions()
		var handlerPath string
		if proto.HasExtension(opts, gen.E_HandlerPath) {
			ext := proto.GetExtension(opts, gen.E_HandlerPath)

			var ok bool
			handlerPath, ok = ext.(string)
			if !ok {
				log.Warn().Msgf("handler path is not a string %s", fullname)
				continue
			}
		}

		sb, err := builder.FromService(svc)
		if err != nil {
			return errors.Wrapf(err, "error building method")
		}

		f, err := builder.FromFile(fd)
		if err != nil {
			return errors.Wrapf(err, "error building proto file")
		}

		// TODO breadchris control if the method should stream or not

		reqMsg := builder.NewMessage(method + "Request")
		reqMsg = reqMsg.AddField(builder.NewField("message", builder.FieldTypeString()))
		req := builder.RpcTypeMessage(reqMsg, false)

		resMsg := builder.NewMessage(method + "Response")
		resMsg = resMsg.AddField(builder.NewField("result", builder.FieldTypeString()))
		res := builder.RpcTypeMessage(resMsg, false)

		b := builder.NewMethod(method, req, res)

		//var mb *builder.MethodBuilder

		// see if this method already exists
		m := svc.FindMethodByName(method)
		if m != nil {
			//// refactor the discovered method to match the new node type
			//mb, err = builder.FromMethod(m)
			//if err != nil {
			//	return nil, errors.Wrapf(err, "error building method")
			//}
			//mb = mb.SetRequestType(inputType).SetResponseType(outputType)
			//sb = sb.RemoveMethod(m.GetName())
			return errors.Errorf("method %s already exists", method)
		}

		// replace the existing method and service with the new one
		sb = sb.AddMethod(b)
		f = f.AddMessage(reqMsg).AddMessage(resMsg)
		f = f.RemoveService(svc.GetName()).AddService(sb)

		if handlerPath != "" {
			_, err := os.Stat(handlerPath)
			if err != nil {
				return errors.Wrapf(err, "error checking handler path: %s", handlerPath)
			}
			handlers = append(handlers, handlerPath)
		}

		nfd, err := f.Build()
		if err != nil {
			return errors.Wrapf(err, "error building proto file")
		}
		err = writeProtoFile(path.Join(dir, protofile), nfd)
		if err != nil {
			return errors.Wrapf(err, "error writing proto file")
		}
	}

	// TODO breadchris this will only work for one handler because of live reload
	for _, handler := range handlers {
		err = addGoMethodToService(handler, method)
		if err != nil {
			return errors.Wrapf(err, "error adding go method to service")
		}
	}

	cmd := exec.Command("go", "generate")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error running go generate: %v, output: %s", err, output)
	}
	return nil
}
