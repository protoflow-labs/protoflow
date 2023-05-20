package grpc

import (
	"context"
	"github.com/protoflow-labs/protoflow/pkg/grpc/bufcurl"
	"testing"
)

func TestExecuteCurl(t *testing.T) {
	input := bufcurl.NewMemoryInputStream()

	outputStream := bufcurl.NewMemoryOutputStream()

	go func() {
		err := ExecuteCurl(context.Background(), RemoteMethod{
			OutputStream:          outputStream,
			TLSConfig:             bufcurl.TLSSettings{},
			URL:                   "http://localhost:8086/protoflow.nodeService/CustomA",
			Protocol:              "grpc",
			Headers:               nil,
			UserAgent:             "",
			ReflectProtocol:       "grpc-v1",
			ReflectHeaders:        nil,
			UnixSocket:            "",
			HTTP2PriorKnowledge:   true,
			NoKeepAlive:           false,
			KeepAliveTimeSeconds:  0,
			ConnectTimeoutSeconds: 0,
		}, input)
		if err != nil {
			t.Error(err)
		}
		outputStream.Close()
	}()
	type TestInput struct {
		Value string
	}
	input.Push(TestInput{Value: "test"})
	input.Close()
	for {
		output, err := outputStream.Next()
		if err != nil {
			break
		}
		t.Log(output)
	}
}
