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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/pkg/grpc/protoencoding"
	"github.com/protoflow-labs/protoflow/pkg/grpc/verbose"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"

	"github.com/bufbuild/connect-go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

type deferredMessage struct {
	data []byte
}

type protoCodec struct{}

func (p protoCodec) Name() string {
	return "proto"
}

func (p protoCodec) Marshal(a any) ([]byte, error) {
	protoMessage, ok := a.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("cannot marshal: %T does not implement proto.Message", a)
	}
	return protoencoding.NewWireMarshaler().Marshal(protoMessage)
}

func (p protoCodec) Unmarshal(bytes []byte, a any) error {
	if deferred, ok := a.(*deferredMessage); ok {
		// must make a copy since Connect framework will re-use the byte slice
		deferred.data = make([]byte, len(bytes))
		copy(deferred.data, bytes)
		return nil
	}
	protoMessage, ok := a.(proto.Message)
	if !ok {
		return fmt.Errorf("cannot unmarshal: %T does not implement proto.Message", a)
	}
	return protoencoding.NewWireUnmarshaler(nil).Unmarshal(bytes, protoMessage)
}

type invokeClient = connect.Client[dynamicpb.Message, deferredMessage]

type invoker struct {
	md           protoreflect.MethodDescriptor
	res          protoencoding.Resolver
	client       *invokeClient
	output       io.Writer
	errOutput    io.Writer
	printer      verbose.Printer
	outputStream OutputStream
}

// NewInvoker creates a new invoker for invoking the method described by the
// given descriptor. The given writer is used to write the output response(s)
// in JSON format. The given resolver is used to resolve Any messages and
// extensions that appear in the input or output. Other parameters are used
// to create a Connect client, for issuing the RPC.
func NewInvoker(container *Container, md protoreflect.MethodDescriptor, res protoencoding.Resolver, httpClient connect.HTTPClient, opts []connect.ClientOption, url string, out io.Writer) Invoker {
	opts = append(opts, connect.WithCodec(protoCodec{}))
	// TODO: could also provide custom compressor implementations that could give us
	//  optics into when request and response messages are compressed (which could be
	//  useful to include in verbose output).
	return &invoker{
		md:           md,
		res:          res,
		output:       out,
		printer:      container.VerbosePrinter(),
		errOutput:    container.Stderr(),
		client:       connect.NewClient[dynamicpb.Message, deferredMessage](httpClient, url, opts...),
		outputStream: NopOutputStream{},
	}
}

func (inv *invoker) InvokeWithStream(ctx context.Context, input InputStream, output OutputStream, headers http.Header) error {
	inv.outputStream = output
	return inv.Invoke(ctx, input, headers)
}

func (inv *invoker) Invoke(ctx context.Context, input InputStream, headers http.Header) error {
	inv.printer.Printf("* Invoking RPC %s\n", inv.md.FullName())
	// request's user-agent header(s) get overwritten by protocol, so we stash them in the
	// context so that underlying transport can restore them
	ctx = withUserAgent(ctx, headers)
	switch {
	case inv.md.IsStreamingServer() && inv.md.IsStreamingClient():
		return inv.handleBidiStream(ctx, input, headers)
	case inv.md.IsStreamingServer():
		return inv.handleServerStream(ctx, input, headers)
	case inv.md.IsStreamingClient():
		return inv.handleClientStream(ctx, input, headers)
	default:
		return inv.handleUnary(ctx, input, headers)
	}
}

func (inv *invoker) handleUnary(ctx context.Context, input InputStream, headers http.Header) error {
	provider := newChanStreamMessageProvider(input, inv.res)
	msg := dynamicpb.NewMessage(inv.md.Input())
	if err := provider.next(msg); err != nil {
		return err
	}
	// make sure input does not contain a second message
	dummy := dynamicpb.NewMessage(inv.md.Input())
	if err := provider.next(dummy); err != io.EOF {
		return fmt.Errorf("method %s is a unary RPC, but input contained more than one request message", inv.md.Name())
	}

	req := connect.NewRequest(msg)
	for k, v := range headers {
		req.Header()[k] = v
	}
	resp, err := inv.client.CallUnary(ctx, req)
	if err != nil {
		var connErr *connect.Error
		if !errors.As(err, &connErr) {
			return err
		}
		err := inv.handleErrorResponse(connErr)
		return err
	}
	err = inv.handleResponse(resp.Msg.data, msg)
	if err != nil {
		inv.outputStream.Error(err)
		return err
	}
	inv.outputStream.Push(msg)
	return nil
}

func (inv *invoker) handleClientStream(ctx context.Context, input InputStream, headers http.Header) (retErr error) {
	provider := newChanStreamMessageProvider(input, inv.res)
	msg := dynamicpb.NewMessage(inv.md.Input())
	stream := inv.client.CallClientStream(ctx)
	for k, v := range headers {
		stream.RequestHeader()[k] = v
	}
	defer func() {
		if retErr != nil {
			var connErr *connect.Error
			if errors.As(retErr, &connErr) {
				retErr = inv.handleErrorResponse(connErr)
			}
		}
	}()
	if err, isStreamError := inv.handleStreamRequest(provider, msg, stream); err != nil {
		if isStreamError {
			_, recvErr := stream.CloseAndReceive()
			// stream.Send should return io.EOF on error, and caller is expected to call
			// stream.Receive to get the actual RPC error.
			if recvErr != nil {
				return recvErr
			}
		}
		return err
	}
	resp, err := stream.CloseAndReceive()
	if err != nil {
		return err
	}

	err = inv.handleResponse(resp.Msg.data, msg)
	if err != nil {
		// TODO breadchris we want to capture all errors, not just the last one
		inv.outputStream.Error(err)
		return err
	}
	inv.outputStream.Push(msg)
	return nil
}

func (inv *invoker) handleServerStream(ctx context.Context, input InputStream, headers http.Header) (retErr error) {
	provider := newChanStreamMessageProvider(input, inv.res)
	msg := dynamicpb.NewMessage(inv.md.Input())
	if err := provider.next(msg); err != nil {
		return err
	}
	// make sure input does not contain a second message
	dummy := dynamicpb.NewMessage(inv.md.Input())
	if err := provider.next(dummy); err != io.EOF {
		return fmt.Errorf("method %s is a unary RPC, but input contained more than one request message", inv.md.Name())
	}

	req := connect.NewRequest(msg)
	for k, v := range headers {
		req.Header()[k] = v
	}
	defer func() {
		if retErr != nil {
			var connErr *connect.Error
			if errors.As(retErr, &connErr) {
				retErr = inv.handleErrorResponse(connErr)
			}
		}
	}()

	stream, err := inv.client.CallServerStream(ctx, req)
	if err != nil {
		return err
	}
	return inv.handleStreamResponse(&serverStreamAdapter{stream: stream})
}

func (inv *invoker) handleBidiStream(ctx context.Context, input InputStream, headers http.Header) (retErr error) {
	ctx, cancel := context.WithCancel(ctx)
	provider := newChanStreamMessageProvider(input, inv.res)
	msg := dynamicpb.NewMessage(inv.md.Input())
	stream := inv.client.CallBidiStream(ctx)
	for k, v := range headers {
		stream.RequestHeader()[k] = v
	}

	defer func() {
		if retErr != nil {
			var connErr *connect.Error
			if errors.As(retErr, &connErr) {
				retErr = inv.handleErrorResponse(connErr)
			}
		}
	}()

	var recvErr error
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer cancel()
		if err := inv.handleStreamResponse(stream); err != nil {
			recvErr = err
		}
	}()
	defer func() {
		wg.Wait()
		if recvErr != nil {
			// may just get io.EOF or cancel error when trying to write to closed
			// request stream whereas actual error details will be seen on the read side
			if retErr == nil || errors.Is(retErr, io.EOF) || isCancelled(retErr) {
				retErr = recvErr
			}
		}
	}()
	shouldCancel := true
	defer func() {
		if shouldCancel {
			cancel()
		}
	}()

	err, isStreamError := inv.handleStreamRequest(provider, msg, stream)
	shouldCancel = err != nil && !isStreamError
	if err != nil {
		return err
	}
	return stream.CloseRequest()
}

func isCancelled(err error) bool {
	if errors.Is(err, context.Canceled) {
		return true
	}
	var connErr *connect.Error
	if errors.As(err, &connErr) {
		return connErr.Code() == connect.CodeCanceled
	}
	return false
}

func (inv *invoker) handleResponse(data []byte, msg *dynamicpb.Message) error {
	if msg == nil {
		msg = dynamicpb.NewMessage(inv.md.Output())
	}
	if err := protoencoding.NewWireUnmarshaler(inv.res).Unmarshal(data, msg); err != nil {
		return err
	}
	// If we want to add a pretty-print option, we should perhaps make this a flag
	// and use protoencoding.NewJSONMarshalerIndent
	outputBytes, err := protoencoding.NewJSONMarshaler(inv.res).Marshal(msg)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(inv.output, "%s\n", outputBytes)
	inv.outputStream.Push(msg)
	return err
}

type clientStream interface {
	Send(message *dynamicpb.Message) error
}

type serverStream interface {
	Receive() (*deferredMessage, error)
	CloseResponse() error
}

type serverStreamAdapter struct {
	stream *connect.ServerStreamForClient[deferredMessage]
}

func (ssa *serverStreamAdapter) Receive() (*deferredMessage, error) {
	if !ssa.stream.Receive() {
		err := ssa.stream.Err()
		if err == nil {
			err = io.EOF
		}
		return nil, err
	}
	return ssa.stream.Msg(), nil
}

func (ssa *serverStreamAdapter) CloseResponse() error {
	return ssa.stream.Close()
}

func (inv *invoker) handleStreamRequest(provider messageProvider, msg *dynamicpb.Message, stream clientStream) (error, bool) {
	for {
		if err := provider.next(msg); errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return err, false
		}
		if err := stream.Send(msg); err != nil {
			return err, true
		}
	}
	return nil, false
}

func (inv *invoker) handleStreamResponse(stream serverStream) error {
	defer func() {
		err := stream.CloseResponse()
		if err != nil {
			inv.outputStream.Error(err)
		}
	}()
	msg := dynamicpb.NewMessage(inv.md.Output())
	for {
		responseMsg, err := stream.Receive()
		if errors.Is(err, io.EOF) {
			inv.outputStream.Error(err)
			return err
		} else if err != nil {
			inv.outputStream.Error(err)
			return err
		}
		if err := inv.handleResponse(responseMsg.data, msg); err != nil {
			inv.outputStream.Error(err)
			return err
		}
		inv.outputStream.Push(msg)
	}
}

func (inv *invoker) handleErrorResponse(connErr *connect.Error) error {
	// NB: This is a nasty hack: we create a fake request that looks
	//     like a unary Connect request, so that the ErrorWriter will
	//     print the error in the format we want, which is just the
	//     JSON representation of the Connect error. (We don't need
	//     an end-stream message representation or for the content
	//     to be put into response headers, which is what it may
	//     choose to do if it detects a different protocol in the
	//     request).
	req := &http.Request{
		Header: http.Header{},
	}
	req.Header.Add("content-type", "application/json")

	w := connect.NewErrorWriter()
	responseWriter := httptest.NewRecorder()
	err := w.Write(responseWriter, req, connErr)
	if err != nil {
		return err
	}
	var prettyPrinted bytes.Buffer
	if err := json.Indent(&prettyPrinted, responseWriter.Body.Bytes(), "", "   "); err != nil {
		return err
	}
	_, _ = inv.errOutput.Write(prettyPrinted.Bytes())
	_, _ = inv.errOutput.Write([]byte("\n"))
	return fmt.Errorf("%d", int(connErr.Code()*8))
}

func newChanStreamMessageProvider(input InputStream, res protoencoding.Resolver) messageProvider {
	return &chanMessageProvider{input: input, res: res}
}

type messageProvider interface {
	next(proto.Message) error
}

type chanMessageProvider struct {
	input InputStream
	res   protoencoding.Resolver
}

func (s *chanMessageProvider) next(msg proto.Message) error {
	data, err := s.input.Next()
	if err != nil {
		return err
	}
	encodedData, err := json.Marshal(data)
	if err != nil {
		return errors.Wrapf(err, "marshaling data")
	}
	proto.Reset(msg)
	return protoencoding.NewJSONUnmarshaler(s.res).Unmarshal(encodedData, msg)
}
