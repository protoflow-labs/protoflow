package grpc

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/fullstorydev/grpcurl"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"sort"
	"strconv"
	"strings"
)

type RpcResult struct {
	DescSource grpcurl.DescriptorSource
	Headers    []rpcMetadata        `json:"headers"`
	Error      *rpcError            `json:"error"`
	Responses  []rpcResponseElement `json:"responses"`
	Requests   *rpcRequestStats     `json:"requests"`
	Trailers   []rpcMetadata        `json:"trailers"`
}

type rpcMetadata struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type rpcInput struct {
	TimeoutSeconds float32           `json:"timeout_seconds"`
	Metadata       []rpcMetadata     `json:"metadata"`
	Data           []json.RawMessage `json:"data"`
}

type rpcResponseElement struct {
	Data    json.RawMessage `json:"message"`
	IsError bool            `json:"isError"`
}

type rpcRequestStats struct {
	Total int `json:"total"`
	Sent  int `json:"sent"`
}

type rpcError struct {
	Code    uint32               `json:"code"`
	Name    string               `json:"name"`
	Message string               `json:"message"`
	Details []rpcResponseElement `json:"details"`
}

type rpcResult struct {
	descSource grpcurl.DescriptorSource
	Headers    []rpcMetadata        `json:"headers"`
	Error      *rpcError            `json:"error"`
	Responses  []rpcResponseElement `json:"responses"`
	Requests   *rpcRequestStats     `json:"requests"`
	Trailers   []rpcMetadata        `json:"trailers"`
}

func (*RpcResult) OnResolveMethod(*desc.MethodDescriptor) {}

func (*RpcResult) OnSendHeaders(metadata.MD) {}

func (r *RpcResult) OnReceiveHeaders(md metadata.MD) {
	r.Headers = responseMetadata(md)
}

func (r *RpcResult) OnReceiveResponse(m proto.Message) {
	r.Responses = append(r.Responses, responseToJSON(r.DescSource, m))
}

func (r *RpcResult) OnReceiveTrailers(stat *status.Status, md metadata.MD) {
	r.Trailers = responseMetadata(md)
	r.Error = toRpcError(r.DescSource, stat)
}

func responseMetadata(md metadata.MD) []rpcMetadata {
	keys := make([]string, 0, len(md))
	for k := range md {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	ret := make([]rpcMetadata, 0, len(md))
	for _, k := range keys {
		vals := md[k]
		for _, v := range vals {
			if strings.HasSuffix(k, "-bin") {
				v = base64.StdEncoding.EncodeToString([]byte(v))
			}
			ret = append(ret, rpcMetadata{Name: k, Value: v})
		}
	}
	return ret
}

func toRpcError(descSource grpcurl.DescriptorSource, stat *status.Status) *rpcError {
	if stat.Code() == codes.OK {
		return nil
	}

	details := stat.Proto().Details
	msgs := make([]rpcResponseElement, len(details))
	for i, d := range details {
		msgs[i] = responseToJSON(descSource, d)
	}
	return &rpcError{
		Code:    uint32(stat.Code()),
		Name:    stat.Code().String(),
		Message: stat.Message(),
		Details: msgs,
	}
}

func responseToJSON(descSource grpcurl.DescriptorSource, msg proto.Message) rpcResponseElement {
	anyResolver := grpcurl.AnyResolverFromDescriptorSourceWithFallback(descSource)
	jsm := jsonpb.Marshaler{EmitDefaults: true, OrigName: true, Indent: "  ", AnyResolver: anyResolver}
	var b bytes.Buffer
	if err := jsm.Marshal(&b, msg); err == nil {
		return rpcResponseElement{Data: json.RawMessage(b.Bytes())}
	} else {
		b, err := json.Marshal(err.Error())
		if err != nil {
			// unable to marshal err message to JSON?
			// should never happen... here's a dumb fallback
			b = []byte(strconv.Quote(err.Error()))
		}
		return rpcResponseElement{Data: b, IsError: true}
	}
}
