package bufcurl

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/pkg/grpc/bufcurl/protoencoding"
	"google.golang.org/protobuf/proto"
)

func resolveMsgFromData(res protoencoding.Resolver, msg proto.Message, data any) error {
	encodedData, err := json.Marshal(data)
	if err != nil {
		return errors.Wrapf(err, "marshaling data")
	}
	proto.Reset(msg)
	return protoencoding.NewJSONUnmarshaler(res).Unmarshal(encodedData, msg)
}
