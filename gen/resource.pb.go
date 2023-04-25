// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        (unknown)
// source: resource.proto

package gen

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

type Resource struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// Types that are assignable to Type:
	//
	//	*Resource_GrpcService
	//	*Resource_RestService
	//	*Resource_Database
	Type   isResource_Type `protobuf_oneof:"type"`
	Blocks []*Block        `protobuf:"bytes,7,rep,name=blocks,proto3" json:"blocks,omitempty"`
}

func (x *Resource) Reset() {
	*x = Resource{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resource_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Resource) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Resource) ProtoMessage() {}

func (x *Resource) ProtoReflect() protoreflect.Message {
	mi := &file_resource_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Resource.ProtoReflect.Descriptor instead.
func (*Resource) Descriptor() ([]byte, []int) {
	return file_resource_proto_rawDescGZIP(), []int{0}
}

func (x *Resource) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Resource) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (m *Resource) GetType() isResource_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (x *Resource) GetGrpcService() *GRPCService {
	if x, ok := x.GetType().(*Resource_GrpcService); ok {
		return x.GrpcService
	}
	return nil
}

func (x *Resource) GetRestService() *RESTService {
	if x, ok := x.GetType().(*Resource_RestService); ok {
		return x.RestService
	}
	return nil
}

func (x *Resource) GetDatabase() *Database {
	if x, ok := x.GetType().(*Resource_Database); ok {
		return x.Database
	}
	return nil
}

func (x *Resource) GetBlocks() []*Block {
	if x != nil {
		return x.Blocks
	}
	return nil
}

type isResource_Type interface {
	isResource_Type()
}

type Resource_GrpcService struct {
	GrpcService *GRPCService `protobuf:"bytes,3,opt,name=grpc_service,json=grpcService,proto3,oneof"`
}

type Resource_RestService struct {
	RestService *RESTService `protobuf:"bytes,4,opt,name=rest_service,json=restService,proto3,oneof"`
}

type Resource_Database struct {
	Database *Database `protobuf:"bytes,5,opt,name=database,proto3,oneof"`
}

func (*Resource_GrpcService) isResource_Type() {}

func (*Resource_RestService) isResource_Type() {}

func (*Resource_Database) isResource_Type() {}

type GRPCService struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Host string `protobuf:"bytes,1,opt,name=host,proto3" json:"host,omitempty"`
}

func (x *GRPCService) Reset() {
	*x = GRPCService{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resource_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GRPCService) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GRPCService) ProtoMessage() {}

func (x *GRPCService) ProtoReflect() protoreflect.Message {
	mi := &file_resource_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GRPCService.ProtoReflect.Descriptor instead.
func (*GRPCService) Descriptor() ([]byte, []int) {
	return file_resource_proto_rawDescGZIP(), []int{1}
}

func (x *GRPCService) GetHost() string {
	if x != nil {
		return x.Host
	}
	return ""
}

type RESTService struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Host   string `protobuf:"bytes,1,opt,name=host,proto3" json:"host,omitempty"`
	Schema string `protobuf:"bytes,2,opt,name=schema,proto3" json:"schema,omitempty"`
}

func (x *RESTService) Reset() {
	*x = RESTService{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resource_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RESTService) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RESTService) ProtoMessage() {}

func (x *RESTService) ProtoReflect() protoreflect.Message {
	mi := &file_resource_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RESTService.ProtoReflect.Descriptor instead.
func (*RESTService) Descriptor() ([]byte, []int) {
	return file_resource_proto_rawDescGZIP(), []int{2}
}

func (x *RESTService) GetHost() string {
	if x != nil {
		return x.Host
	}
	return ""
}

func (x *RESTService) GetSchema() string {
	if x != nil {
		return x.Schema
	}
	return ""
}

type Database struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Host string `protobuf:"bytes,1,opt,name=host,proto3" json:"host,omitempty"`
}

func (x *Database) Reset() {
	*x = Database{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resource_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Database) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Database) ProtoMessage() {}

func (x *Database) ProtoReflect() protoreflect.Message {
	mi := &file_resource_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Database.ProtoReflect.Descriptor instead.
func (*Database) Descriptor() ([]byte, []int) {
	return file_resource_proto_rawDescGZIP(), []int{3}
}

func (x *Database) GetHost() string {
	if x != nil {
		return x.Host
	}
	return ""
}

var File_resource_proto protoreflect.FileDescriptor

var file_resource_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x08, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x1a, 0x0b, 0x62, 0x6c, 0x6f, 0x63,
	0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x86, 0x02, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x3a, 0x0a, 0x0c, 0x67, 0x72, 0x70, 0x63,
	0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15,
	0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x47, 0x52, 0x50, 0x43, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x48, 0x00, 0x52, 0x0b, 0x67, 0x72, 0x70, 0x63, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x12, 0x3a, 0x0a, 0x0c, 0x72, 0x65, 0x73, 0x74, 0x5f, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x52, 0x45, 0x53, 0x54, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x48, 0x00, 0x52, 0x0b, 0x72, 0x65, 0x73, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x30, 0x0a, 0x08, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x12, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x44, 0x61,
	0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x48, 0x00, 0x52, 0x08, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61,
	0x73, 0x65, 0x12, 0x24, 0x0a, 0x06, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x18, 0x07, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x2e, 0x42, 0x6c, 0x6f, 0x63, 0x6b,
	0x52, 0x06, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x42, 0x06, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x22, 0x21, 0x0a, 0x0b, 0x47, 0x52, 0x50, 0x43, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x12, 0x0a, 0x04, 0x68, 0x6f, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68,
	0x6f, 0x73, 0x74, 0x22, 0x39, 0x0a, 0x0b, 0x52, 0x45, 0x53, 0x54, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x6f, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x68, 0x6f, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x22, 0x1e,
	0x0a, 0x08, 0x44, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x6f,
	0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x6f, 0x73, 0x74, 0x42, 0x86,
	0x01, 0x0a, 0x0c, 0x63, 0x6f, 0x6d, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x42,
	0x0d, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01,
	0x5a, 0x27, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x66, 0x6c, 0x6f, 0x77, 0x2d, 0x6c, 0x61, 0x62, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x66, 0x6c, 0x6f, 0x77, 0x2f, 0x67, 0x65, 0x6e, 0xa2, 0x02, 0x03, 0x52, 0x58, 0x58, 0xaa,
	0x02, 0x08, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0xca, 0x02, 0x08, 0x52, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0xe2, 0x02, 0x14, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x08, 0x52,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_resource_proto_rawDescOnce sync.Once
	file_resource_proto_rawDescData = file_resource_proto_rawDesc
)

func file_resource_proto_rawDescGZIP() []byte {
	file_resource_proto_rawDescOnce.Do(func() {
		file_resource_proto_rawDescData = protoimpl.X.CompressGZIP(file_resource_proto_rawDescData)
	})
	return file_resource_proto_rawDescData
}

var file_resource_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_resource_proto_goTypes = []interface{}{
	(*Resource)(nil),    // 0: resource.Resource
	(*GRPCService)(nil), // 1: resource.GRPCService
	(*RESTService)(nil), // 2: resource.RESTService
	(*Database)(nil),    // 3: resource.Database
	(*Block)(nil),       // 4: block.Block
}
var file_resource_proto_depIdxs = []int32{
	1, // 0: resource.Resource.grpc_service:type_name -> resource.GRPCService
	2, // 1: resource.Resource.rest_service:type_name -> resource.RESTService
	3, // 2: resource.Resource.database:type_name -> resource.Database
	4, // 3: resource.Resource.blocks:type_name -> block.Block
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_resource_proto_init() }
func file_resource_proto_init() {
	if File_resource_proto != nil {
		return
	}
	file_block_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_resource_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Resource); i {
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
		file_resource_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GRPCService); i {
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
		file_resource_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RESTService); i {
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
		file_resource_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Database); i {
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
	file_resource_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*Resource_GrpcService)(nil),
		(*Resource_RestService)(nil),
		(*Resource_Database)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_resource_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_resource_proto_goTypes,
		DependencyIndexes: file_resource_proto_depIdxs,
		MessageInfos:      file_resource_proto_msgTypes,
	}.Build()
	File_resource_proto = out.File
	file_resource_proto_rawDesc = nil
	file_resource_proto_goTypes = nil
	file_resource_proto_depIdxs = nil
}
