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

package bufcurl

import (
	"context"
	"fmt"
	"github.com/protoflow-labs/protoflow/pkg/grpc/bufcurl/protoencoding"
	"io"
	"net/http"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

type InputStream interface {
	Push(any)
	Next() (any, error)
	Close()
}

type MemoryInputStream struct {
	Stream chan any
}

var _ InputStream = &MemoryInputStream{}

type MemoryStreamOptions func(*MemoryInputStream)

func NewMemoryInputStream() *MemoryInputStream {
	return &MemoryInputStream{
		Stream: make(chan any),
	}
}

func (m *MemoryInputStream) Push(a any) {
	m.Stream <- a
}

func (m *MemoryInputStream) Next() (any, error) {
	data, ok := <-m.Stream
	if !ok {
		return nil, io.EOF
	}
	return data, nil
}

func (m *MemoryInputStream) Close() {
	close(m.Stream)
}

type OutputStream interface {
	Push(any)
	Error(error)
	Next() (any, error)
	Close()
}

type NopOutputStream struct{}

func (n NopOutputStream) Push(a any) {}

func (n NopOutputStream) Error(err error) {}

func (n NopOutputStream) Next() (any, error) {
	return nil, nil
}

func (n NopOutputStream) Close() {}

var _ OutputStream = &NopOutputStream{}

type MemoryOutputStream struct {
	Out chan any
	Err chan error
}

var _ OutputStream = &MemoryOutputStream{}

func NewMemoryOutputStream() *MemoryOutputStream {
	return &MemoryOutputStream{
		Out: make(chan any),
		Err: make(chan error),
	}
}

func (m *MemoryOutputStream) Push(a any) {
	m.Out <- a
}

func (m *MemoryOutputStream) Error(err error) {
	m.Err <- err
}

func (m *MemoryOutputStream) Next() (any, error) {
	select {
	case data, ok := <-m.Out:
		if !ok {
			return nil, io.EOF
		}
		return data, nil
	case err, ok := <-m.Err:
		if !ok {
			return nil, io.EOF
		}
		return nil, err
	}
}

func (m *MemoryOutputStream) Close() {
	close(m.Out)
	close(m.Err)
}

// Invoker provides the ability to invoke RPCs dynamically.
type Invoker interface {
	// Invoke invokes an RPC method using the given input data and request headers.
	// The dataSource is a string that describes the input data (e.g. a filename).
	// The actual contents of the request data is read from the given reader.
	Invoke(ctx context.Context, input InputStream, headers http.Header) error
	InvokeWithStream(ctx context.Context, input InputStream, output OutputStream, headers http.Header) error
}

// ResolveMethodDescriptor uses the given resolver to find a descriptor for
// the requested service and method. The service name must be fully-qualified.
func ResolveMethodDescriptor(res protoencoding.Resolver, service, method string) (protoreflect.MethodDescriptor, error) {
	descriptor, err := res.FindDescriptorByName(protoreflect.FullName(service))
	if err == protoregistry.NotFound {
		return nil, fmt.Errorf("failed to find service named %q in schema", service)
	} else if err != nil {
		return nil, err
	}
	serviceDescriptor, ok := descriptor.(protoreflect.ServiceDescriptor)
	if !ok {
		return nil, fmt.Errorf("URL indicates service name %q, but that name is a %s", service, DescriptorKind(descriptor))
	}
	methodDescriptor := serviceDescriptor.Methods().ByName(protoreflect.Name(method))
	if methodDescriptor == nil {
		return nil, fmt.Errorf("URL indicates method name %q, but service %q contains no such method", method, service)
	}
	return methodDescriptor, nil
}
