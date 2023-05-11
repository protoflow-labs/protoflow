package grpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/fullstorydev/grpcurl"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"io"
)

func CallMethod(conn *grpc.ClientConn, input interface{}, fullQlName string) (interface{}, error) {
	var headers []string
	methods, err := AllMethodsViaReflection(context.Background(), conn)
	if err != nil {
		return nil, errors.Wrap(err, "error getting all methods")
	}

	requestFunc := func(m proto.Message) error {
		req, err := json.Marshal(input)
		if err != nil {
			return errors.Wrapf(err, "error marshalling input: %v", input)
		}

		if err := jsonpb.Unmarshal(bytes.NewReader(req), m); err != nil {
			return errors.Wrapf(err, "error unmarshalling input: %s", string(req))
		}
		return io.EOF
	}
	// TODO breadchris do we have to do this to make this call? There is probably a better way
	for _, m := range methods {
		if m.GetFullyQualifiedName() == fullQlName {
			descSource, err := grpcurl.DescriptorSourceFromFileDescriptors(m.GetFile())
			if err != nil {
				return nil, errors.Wrap(err, "error getting descriptor source")
			}
			result := RpcResult{
				DescSource: descSource,
			}
			if err := grpcurl.InvokeRPC(context.Background(), descSource, conn, fullQlName, headers, &result, requestFunc); err != nil {
				return nil, errors.Wrapf(err, "error invoking rpc %s", fullQlName)
			}

			if len(result.Responses) == 0 {
				return nil, fmt.Errorf("no responses received")
			}

			resp := result.Responses[0]
			var data interface{}
			err = json.Unmarshal(resp.Data, &data)
			if err != nil {
				return nil, errors.Wrapf(err, "error unmarshalling response: %s", string(resp.Data))
			}

			return data, err
		}
	}
	return nil, fmt.Errorf("method not found: %s", fullQlName)
}
