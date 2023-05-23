package manager

import (
	"fmt"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
	"github.com/jhump/protoreflect/desc/protoprint"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func GetProtoForMethod(packageName, serviceName string, method protoreflect.MethodDescriptor) (string, error) {
	md, err := desc.WrapMethod(method)
	if err != nil {
		return "", errors.Wrapf(err, "error wrapping method")
	}

	b, err := builder.FromMethod(md)
	if err != nil {
		return "", errors.Wrapf(err, "error building method descriptor")
	}

	methodType, err := printBuilder(b)
	if err != nil {
		return "", err
	}

	inputType, err := printMessage(method.Input())
	if err != nil {
		return "", errors.Wrapf(err, "error printing input message")
	}
	outputType, err := printMessage(method.Output())
	if err != nil {
		return "", errors.Wrapf(err, "error printing output message")
	}

	methodStr := fmt.Sprintf("package %s;\n", packageName)
	methodStr += fmt.Sprintf("service %s {\n", serviceName)
	methodStr += "\t" + methodType
	methodStr += "}\n"
	methodStr += inputType
	methodStr += outputType
	return methodStr, nil
}

func printBuilder(b builder.Builder) (string, error) {
	d, err := b.BuildDescriptor()
	if err != nil {
		return "", errors.Wrapf(err, "error building method descriptor")
	}

	p := protoprint.Printer{}
	s, err := p.PrintProtoToString(d)
	if err != nil {
		return "", errors.Wrapf(err, "error printing proto")
	}
	return s, nil
}

func printMessage(msgType protoreflect.MessageDescriptor) (string, error) {
	msg, err := desc.WrapMessage(msgType)
	if err != nil {
		return "", errors.Wrapf(err, "error wrapping message")
	}

	m, err := builder.FromMessage(msg)
	if err != nil {
		return "", errors.Wrapf(err, "error building message descriptor")
	}
	return printBuilder(m)
}
