// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
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

type Runtime int32

const (
	Runtime_NODEJS Runtime = 0
	Runtime_PYTHON Runtime = 1
	Runtime_GO     Runtime = 2
)

// Enum value maps for Runtime.
var (
	Runtime_name = map[int32]string{
		0: "NODEJS",
		1: "PYTHON",
		2: "GO",
	}
	Runtime_value = map[string]int32{
		"NODEJS": 0,
		"PYTHON": 1,
		"GO":     2,
	}
)

func (x Runtime) Enum() *Runtime {
	p := new(Runtime)
	*p = x
	return p
}

func (x Runtime) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Runtime) Descriptor() protoreflect.EnumDescriptor {
	return file_resource_proto_enumTypes[0].Descriptor()
}

func (Runtime) Type() protoreflect.EnumType {
	return &file_resource_proto_enumTypes[0]
}

func (x Runtime) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Runtime.Descriptor instead.
func (Runtime) EnumDescriptor() ([]byte, []int) {
	return file_resource_proto_rawDescGZIP(), []int{0}
}

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
	//	*Resource_Docstore
	//	*Resource_Blobstore
	//	*Resource_LanguageService
	Type isResource_Type `protobuf_oneof:"type"`
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

func (x *Resource) GetDocstore() *Docstore {
	if x, ok := x.GetType().(*Resource_Docstore); ok {
		return x.Docstore
	}
	return nil
}

func (x *Resource) GetBlobstore() *Blobstore {
	if x, ok := x.GetType().(*Resource_Blobstore); ok {
		return x.Blobstore
	}
	return nil
}

func (x *Resource) GetLanguageService() *LanguageService {
	if x, ok := x.GetType().(*Resource_LanguageService); ok {
		return x.LanguageService
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

type Resource_Docstore struct {
	Docstore *Docstore `protobuf:"bytes,5,opt,name=docstore,proto3,oneof"`
}

type Resource_Blobstore struct {
	Blobstore *Blobstore `protobuf:"bytes,6,opt,name=blobstore,proto3,oneof"`
}

type Resource_LanguageService struct {
	LanguageService *LanguageService `protobuf:"bytes,7,opt,name=language_service,json=languageService,proto3,oneof"`
}

func (*Resource_GrpcService) isResource_Type() {}

func (*Resource_RestService) isResource_Type() {}

func (*Resource_Docstore) isResource_Type() {}

func (*Resource_Blobstore) isResource_Type() {}

func (*Resource_LanguageService) isResource_Type() {}

type LanguageService struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Runtime Runtime      `protobuf:"varint,1,opt,name=runtime,proto3,enum=resource.Runtime" json:"runtime,omitempty"`
	Grpc    *GRPCService `protobuf:"bytes,2,opt,name=grpc,proto3" json:"grpc,omitempty"`
}

func (x *LanguageService) Reset() {
	*x = LanguageService{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resource_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LanguageService) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LanguageService) ProtoMessage() {}

func (x *LanguageService) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use LanguageService.ProtoReflect.Descriptor instead.
func (*LanguageService) Descriptor() ([]byte, []int) {
	return file_resource_proto_rawDescGZIP(), []int{1}
}

func (x *LanguageService) GetRuntime() Runtime {
	if x != nil {
		return x.Runtime
	}
	return Runtime_NODEJS
}

func (x *LanguageService) GetGrpc() *GRPCService {
	if x != nil {
		return x.Grpc
	}
	return nil
}

type GRPCService struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Host string `protobuf:"bytes,1,opt,name=host,proto3" json:"host,omitempty"`
}

func (x *GRPCService) Reset() {
	*x = GRPCService{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resource_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GRPCService) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GRPCService) ProtoMessage() {}

func (x *GRPCService) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use GRPCService.ProtoReflect.Descriptor instead.
func (*GRPCService) Descriptor() ([]byte, []int) {
	return file_resource_proto_rawDescGZIP(), []int{2}
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

	BaseUrl string `protobuf:"bytes,1,opt,name=base_url,json=baseUrl,proto3" json:"base_url,omitempty"`
}

func (x *RESTService) Reset() {
	*x = RESTService{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resource_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RESTService) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RESTService) ProtoMessage() {}

func (x *RESTService) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use RESTService.ProtoReflect.Descriptor instead.
func (*RESTService) Descriptor() ([]byte, []int) {
	return file_resource_proto_rawDescGZIP(), []int{3}
}

func (x *RESTService) GetBaseUrl() string {
	if x != nil {
		return x.BaseUrl
	}
	return ""
}

type Docstore struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
}

func (x *Docstore) Reset() {
	*x = Docstore{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resource_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Docstore) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Docstore) ProtoMessage() {}

func (x *Docstore) ProtoReflect() protoreflect.Message {
	mi := &file_resource_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Docstore.ProtoReflect.Descriptor instead.
func (*Docstore) Descriptor() ([]byte, []int) {
	return file_resource_proto_rawDescGZIP(), []int{4}
}

func (x *Docstore) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

type Blobstore struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
}

func (x *Blobstore) Reset() {
	*x = Blobstore{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resource_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Blobstore) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Blobstore) ProtoMessage() {}

func (x *Blobstore) ProtoReflect() protoreflect.Message {
	mi := &file_resource_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Blobstore.ProtoReflect.Descriptor instead.
func (*Blobstore) Descriptor() ([]byte, []int) {
	return file_resource_proto_rawDescGZIP(), []int{5}
}

func (x *Blobstore) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

var File_resource_proto protoreflect.FileDescriptor

var file_resource_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x08, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x1a, 0x0b, 0x62, 0x6c, 0x6f, 0x63,
	0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xdd, 0x02, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x6f,
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
	0x12, 0x30, 0x0a, 0x08, 0x64, 0x6f, 0x63, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x12, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x44, 0x6f,
	0x63, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x48, 0x00, 0x52, 0x08, 0x64, 0x6f, 0x63, 0x73, 0x74, 0x6f,
	0x72, 0x65, 0x12, 0x33, 0x0a, 0x09, 0x62, 0x6c, 0x6f, 0x62, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x2e, 0x42, 0x6c, 0x6f, 0x62, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x48, 0x00, 0x52, 0x09, 0x62, 0x6c,
	0x6f, 0x62, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x12, 0x46, 0x0a, 0x10, 0x6c, 0x61, 0x6e, 0x67, 0x75,
	0x61, 0x67, 0x65, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x19, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x4c, 0x61, 0x6e,
	0x67, 0x75, 0x61, 0x67, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x48, 0x00, 0x52, 0x0f,
	0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x42,
	0x06, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x69, 0x0a, 0x0f, 0x4c, 0x61, 0x6e, 0x67, 0x75,
	0x61, 0x67, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x2b, 0x0a, 0x07, 0x72, 0x75,
	0x6e, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x11, 0x2e, 0x72, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x52, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x52, 0x07,
	0x72, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x29, 0x0a, 0x04, 0x67, 0x72, 0x70, 0x63, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x2e, 0x47, 0x52, 0x50, 0x43, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x52, 0x04, 0x67, 0x72,
	0x70, 0x63, 0x22, 0x21, 0x0a, 0x0b, 0x47, 0x52, 0x50, 0x43, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x6f, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x68, 0x6f, 0x73, 0x74, 0x22, 0x28, 0x0a, 0x0b, 0x52, 0x45, 0x53, 0x54, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x62, 0x61, 0x73, 0x65, 0x5f, 0x75, 0x72, 0x6c,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x62, 0x61, 0x73, 0x65, 0x55, 0x72, 0x6c, 0x22,
	0x1c, 0x0a, 0x08, 0x44, 0x6f, 0x63, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x75,
	0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x22, 0x1d, 0x0a,
	0x09, 0x42, 0x6c, 0x6f, 0x62, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72,
	0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x2a, 0x29, 0x0a, 0x07,
	0x52, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x0a, 0x0a, 0x06, 0x4e, 0x4f, 0x44, 0x45, 0x4a,
	0x53, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x50, 0x59, 0x54, 0x48, 0x4f, 0x4e, 0x10, 0x01, 0x12,
	0x06, 0x0a, 0x02, 0x47, 0x4f, 0x10, 0x02, 0x42, 0x86, 0x01, 0x0a, 0x0c, 0x63, 0x6f, 0x6d, 0x2e,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x42, 0x0d, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x27, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x66, 0x6c, 0x6f, 0x77, 0x2d,
	0x6c, 0x61, 0x62, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x66, 0x6c, 0x6f, 0x77, 0x2f, 0x67,
	0x65, 0x6e, 0xa2, 0x02, 0x03, 0x52, 0x58, 0x58, 0xaa, 0x02, 0x08, 0x52, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0xca, 0x02, 0x08, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0xe2, 0x02,
	0x14, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x08, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
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

var file_resource_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_resource_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_resource_proto_goTypes = []interface{}{
	(Runtime)(0),            // 0: resource.Runtime
	(*Resource)(nil),        // 1: resource.Resource
	(*LanguageService)(nil), // 2: resource.LanguageService
	(*GRPCService)(nil),     // 3: resource.GRPCService
	(*RESTService)(nil),     // 4: resource.RESTService
	(*Docstore)(nil),        // 5: resource.Docstore
	(*Blobstore)(nil),       // 6: resource.Blobstore
}
var file_resource_proto_depIdxs = []int32{
	3, // 0: resource.Resource.grpc_service:type_name -> resource.GRPCService
	4, // 1: resource.Resource.rest_service:type_name -> resource.RESTService
	5, // 2: resource.Resource.docstore:type_name -> resource.Docstore
	6, // 3: resource.Resource.blobstore:type_name -> resource.Blobstore
	2, // 4: resource.Resource.language_service:type_name -> resource.LanguageService
	0, // 5: resource.LanguageService.runtime:type_name -> resource.Runtime
	3, // 6: resource.LanguageService.grpc:type_name -> resource.GRPCService
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
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
			switch v := v.(*LanguageService); i {
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
		file_resource_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
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
		file_resource_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Docstore); i {
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
		file_resource_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Blobstore); i {
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
		(*Resource_Docstore)(nil),
		(*Resource_Blobstore)(nil),
		(*Resource_LanguageService)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_resource_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_resource_proto_goTypes,
		DependencyIndexes: file_resource_proto_depIdxs,
		EnumInfos:         file_resource_proto_enumTypes,
		MessageInfos:      file_resource_proto_msgTypes,
	}.Build()
	File_resource_proto = out.File
	file_resource_proto_rawDesc = nil
	file_resource_proto_goTypes = nil
	file_resource_proto_depIdxs = nil
}
