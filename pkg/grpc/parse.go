package grpc

import (
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/proto"
	"io"
	"os"
	"path"
	"strings"
)

func FileContentsFromMap(files map[string]string) protoparse.FileAccessor {
	return func(filename string) (io.ReadCloser, error) {
		contents, ok := files[filename]
		if !ok {
			return nil, os.ErrNotExist
		}
		return io.NopCloser(strings.NewReader(contents)), nil
	}
}

func ParseProto(file string) ([]*desc.FileDescriptor, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, errors.Wrapf(err, "error reading file %s", file)
	}
	// split filename from path
	_, filename := path.Split(file)

	protoMap := map[string]string{filename: string(content)}

	protoFiles, err := proto.Proto.ReadDir(".")
	if err != nil {
		return nil, errors.Wrapf(err, "error reading proto directory")
	}
	for _, protoFile := range protoFiles {
		if strings.HasSuffix(protoFile.Name(), ".proto") {
			content, err := proto.Proto.ReadFile(protoFile.Name())
			if err != nil {
				return nil, errors.Wrapf(err, "error reading file %s", protoFile.Name())
			}
			protoMap[protoFile.Name()] = string(content)
		}
	}

	parser := protoparse.Parser{
		ImportPaths:                     nil,
		InferImportPaths:                false,
		LookupImport:                    nil,
		LookupImportProto:               nil,
		Accessor:                        FileContentsFromMap(protoMap),
		IncludeSourceCodeInfo:           false,
		ValidateUnlinkedFiles:           false,
		InterpretOptionsInUnlinkedFiles: false,
		ErrorReporter:                   nil,
		WarningReporter:                 nil,
	}
	return parser.ParseFiles(filename)
}
