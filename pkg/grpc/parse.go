package grpc

import (
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/proto"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
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

func ParseProtoDir(dir, protofile string) ([]*desc.FileDescriptor, error) {
	protoMap, err := protoflowProtoMap(dir, os.ReadDir, os.ReadFile)
	if err != nil {
		return nil, errors.Wrapf(err, "error reading proto directory")
	}
	parser := protoparse.Parser{
		ImportPaths:       nil,
		InferImportPaths:  false,
		LookupImport:      nil,
		LookupImportProto: nil,
		// TODO breadchris the accessor should be able to read files, not just returning static file contents
		Accessor:                        FileContentsFromMap(protoMap),
		IncludeSourceCodeInfo:           false,
		ValidateUnlinkedFiles:           false,
		InterpretOptionsInUnlinkedFiles: false,
		ErrorReporter:                   nil,
		WarningReporter:                 nil,
	}
	// TODO breadchris these files should be read from proto dir
	return parser.ParseFiles(protofile)
}

func protoflowProtoMap(
	dir string,
	readDir func(dir string) ([]fs.DirEntry, error),
	readFile func(name string) ([]byte, error),
) (map[string]string, error) {
	protoMap := map[string]string{}

	var processFiles func(string) error
	processFiles = func(p string) error {
		protoFiles, err := readDir(p)
		if err != nil {
			return errors.Wrapf(err, "error reading proto directory")
		}

		for _, protoFile := range protoFiles {
			fullPath := filepath.Join(p, protoFile.Name())

			if protoFile.IsDir() {
				if err := processFiles(p); err != nil {
					return err
				}
			} else if strings.HasSuffix(protoFile.Name(), ".proto") {
				content, err := readFile(fullPath)
				if err != nil {
					return errors.Wrapf(err, "error reading file %s", fullPath)
				}
				protoMap[fullPath] = string(content)
			}
		}
		return nil
	}

	err := processFiles(dir)
	if err != nil {
		return nil, err
	}

	return protoMap, nil
}

func ParseProto(file string) ([]*desc.FileDescriptor, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, errors.Wrapf(err, "error reading file %s", file)
	}
	// split filename from path
	_, filename := path.Split(file)

	protoMap, err := protoflowProtoMap(".", proto.Proto.ReadDir, proto.Proto.ReadFile)
	if err != nil {
		return nil, errors.Wrapf(err, "error reading proto directory")
	}
	protoMap[filename] = string(content)

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
