// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        (unknown)
// source: grpc/grpc.proto

package grpc

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Method struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Package string `protobuf:"bytes,1,opt,name=package,proto3" json:"package,omitempty"`
	Service string `protobuf:"bytes,2,opt,name=service,proto3" json:"service,omitempty"`
	Method  string `protobuf:"bytes,3,opt,name=method,proto3" json:"method,omitempty"`
}

func (x *Method) Reset() {
	*x = Method{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_grpc_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Method) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Method) ProtoMessage() {}

func (x *Method) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_grpc_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Method.ProtoReflect.Descriptor instead.
func (*Method) Descriptor() ([]byte, []int) {
	return file_grpc_grpc_proto_rawDescGZIP(), []int{0}
}

func (x *Method) GetPackage() string {
	if x != nil {
		return x.Package
	}
	return ""
}

func (x *Method) GetService() string {
	if x != nil {
		return x.Service
	}
	return ""
}

func (x *Method) GetMethod() string {
	if x != nil {
		return x.Method
	}
	return ""
}

type Server struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Host string `protobuf:"bytes,1,opt,name=host,proto3" json:"host,omitempty"`
}

func (x *Server) Reset() {
	*x = Server{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_grpc_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Server) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Server) ProtoMessage() {}

func (x *Server) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_grpc_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Server.ProtoReflect.Descriptor instead.
func (*Server) Descriptor() ([]byte, []int) {
	return file_grpc_grpc_proto_rawDescGZIP(), []int{1}
}

func (x *Server) GetHost() string {
	if x != nil {
		return x.Host
	}
	return ""
}

type GRPC struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Type:
	//
	//	*GRPC_Method
	//	*GRPC_Server
	Type isGRPC_Type `protobuf_oneof:"type"`
}

func (x *GRPC) Reset() {
	*x = GRPC{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_grpc_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GRPC) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GRPC) ProtoMessage() {}

func (x *GRPC) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_grpc_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GRPC.ProtoReflect.Descriptor instead.
func (*GRPC) Descriptor() ([]byte, []int) {
	return file_grpc_grpc_proto_rawDescGZIP(), []int{2}
}

func (m *GRPC) GetType() isGRPC_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (x *GRPC) GetMethod() *Method {
	if x, ok := x.GetType().(*GRPC_Method); ok {
		return x.Method
	}
	return nil
}

func (x *GRPC) GetServer() *Server {
	if x, ok := x.GetType().(*GRPC_Server); ok {
		return x.Server
	}
	return nil
}

type isGRPC_Type interface {
	isGRPC_Type()
}

type GRPC_Method struct {
	Method *Method `protobuf:"bytes,1,opt,name=method,proto3,oneof"`
}

type GRPC_Server struct {
	Server *Server `protobuf:"bytes,2,opt,name=server,proto3,oneof"`
}

func (*GRPC_Method) isGRPC_Type() {}

func (*GRPC_Server) isGRPC_Type() {}

var File_grpc_grpc_proto protoreflect.FileDescriptor

var file_grpc_grpc_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x04, 0x67, 0x72, 0x70, 0x63, 0x22, 0x54, 0x0a, 0x06, 0x4d, 0x65, 0x74, 0x68, 0x6f,
	0x64, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x70, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x22, 0x1c, 0x0a,
	0x06, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x6f, 0x73, 0x74, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x6f, 0x73, 0x74, 0x22, 0x5e, 0x0a, 0x04, 0x47,
	0x52, 0x50, 0x43, 0x12, 0x26, 0x0a, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x4d, 0x65, 0x74, 0x68, 0x6f,
	0x64, 0x48, 0x00, 0x52, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x26, 0x0a, 0x06, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x67, 0x72,
	0x70, 0x63, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x48, 0x00, 0x52, 0x06, 0x73, 0x65, 0x72,
	0x76, 0x65, 0x72, 0x42, 0x06, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x42, 0x73, 0x0a, 0x08, 0x63,
	0x6f, 0x6d, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x42, 0x09, 0x47, 0x72, 0x70, 0x63, 0x50, 0x72, 0x6f,
	0x74, 0x6f, 0x50, 0x01, 0x5a, 0x2c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x66, 0x6c, 0x6f, 0x77, 0x2d, 0x6c, 0x61, 0x62, 0x73, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x66, 0x6c, 0x6f, 0x77, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x72,
	0x70, 0x63, 0xa2, 0x02, 0x03, 0x47, 0x58, 0x58, 0xaa, 0x02, 0x04, 0x47, 0x72, 0x70, 0x63, 0xca,
	0x02, 0x04, 0x47, 0x72, 0x70, 0x63, 0xe2, 0x02, 0x10, 0x47, 0x72, 0x70, 0x63, 0x5c, 0x47, 0x50,
	0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x04, 0x47, 0x72, 0x70, 0x63,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_grpc_grpc_proto_rawDescOnce sync.Once
	file_grpc_grpc_proto_rawDescData = file_grpc_grpc_proto_rawDesc
)

func file_grpc_grpc_proto_rawDescGZIP() []byte {
	file_grpc_grpc_proto_rawDescOnce.Do(func() {
		file_grpc_grpc_proto_rawDescData = protoimpl.X.CompressGZIP(file_grpc_grpc_proto_rawDescData)
	})
	return file_grpc_grpc_proto_rawDescData
}

var file_grpc_grpc_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_grpc_grpc_proto_goTypes = []interface{}{
	(*Method)(nil), // 0: grpc.Method
	(*Server)(nil), // 1: grpc.Server
	(*GRPC)(nil),   // 2: grpc.GRPC
}
var file_grpc_grpc_proto_depIdxs = []int32{
	0, // 0: grpc.GRPC.method:type_name -> grpc.Method
	1, // 1: grpc.GRPC.server:type_name -> grpc.Server
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_grpc_grpc_proto_init() }
func file_grpc_grpc_proto_init() {
	if File_grpc_grpc_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_grpc_grpc_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Method); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_grpc_grpc_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Server); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_grpc_grpc_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GRPC); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_grpc_grpc_proto_msgTypes[2].OneofWrappers = []interface{}{
		(*GRPC_Method)(nil),
		(*GRPC_Server)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_grpc_grpc_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_grpc_grpc_proto_goTypes,
		DependencyIndexes: file_grpc_grpc_proto_depIdxs,
		MessageInfos:      file_grpc_grpc_proto_msgTypes,
	}.Build()
	File_grpc_grpc_proto = out.File
	file_grpc_grpc_proto_rawDesc = nil
	file_grpc_grpc_proto_goTypes = nil
	file_grpc_grpc_proto_depIdxs = nil
}
