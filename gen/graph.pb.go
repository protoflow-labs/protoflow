// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        (unknown)
// source: graph.proto

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

type Graph struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    string  `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name  string  `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Nodes []*Node `protobuf:"bytes,3,rep,name=nodes,proto3" json:"nodes,omitempty"`
	Edges []*Edge `protobuf:"bytes,4,rep,name=edges,proto3" json:"edges,omitempty"`
}

func (x *Graph) Reset() {
	*x = Graph{}
	if protoimpl.UnsafeEnabled {
		mi := &file_graph_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Graph) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Graph) ProtoMessage() {}

func (x *Graph) ProtoReflect() protoreflect.Message {
	mi := &file_graph_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Graph.ProtoReflect.Descriptor instead.
func (*Graph) Descriptor() ([]byte, []int) {
	return file_graph_proto_rawDescGZIP(), []int{0}
}

func (x *Graph) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Graph) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Graph) GetNodes() []*Node {
	if x != nil {
		return x.Nodes
	}
	return nil
}

func (x *Graph) GetEdges() []*Edge {
	if x != nil {
		return x.Edges
	}
	return nil
}

type Node struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id      string  `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name    string  `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	BlockId string  `protobuf:"bytes,3,opt,name=block_id,json=blockId,proto3" json:"block_id,omitempty"`
	X       float32 `protobuf:"fixed32,4,opt,name=x,proto3" json:"x,omitempty"`
	Y       float32 `protobuf:"fixed32,5,opt,name=y,proto3" json:"y,omitempty"`
	// Types that are assignable to Config:
	//
	//	*Node_Rest
	//	*Node_Grpc
	//	*Node_Collection
	//	*Node_Entity
	//	*Node_Input
	Config isNode_Config `protobuf_oneof:"config"`
	// Dependencies
	ResourceIds []string `protobuf:"bytes,11,rep,name=resource_ids,json=resourceIds,proto3" json:"resource_ids,omitempty"`
}

func (x *Node) Reset() {
	*x = Node{}
	if protoimpl.UnsafeEnabled {
		mi := &file_graph_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Node) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Node) ProtoMessage() {}

func (x *Node) ProtoReflect() protoreflect.Message {
	mi := &file_graph_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Node.ProtoReflect.Descriptor instead.
func (*Node) Descriptor() ([]byte, []int) {
	return file_graph_proto_rawDescGZIP(), []int{1}
}

func (x *Node) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Node) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Node) GetBlockId() string {
	if x != nil {
		return x.BlockId
	}
	return ""
}

func (x *Node) GetX() float32 {
	if x != nil {
		return x.X
	}
	return 0
}

func (x *Node) GetY() float32 {
	if x != nil {
		return x.Y
	}
	return 0
}

func (m *Node) GetConfig() isNode_Config {
	if m != nil {
		return m.Config
	}
	return nil
}

func (x *Node) GetRest() *REST {
	if x, ok := x.GetConfig().(*Node_Rest); ok {
		return x.Rest
	}
	return nil
}

func (x *Node) GetGrpc() *GRPC {
	if x, ok := x.GetConfig().(*Node_Grpc); ok {
		return x.Grpc
	}
	return nil
}

func (x *Node) GetCollection() *Collection {
	if x, ok := x.GetConfig().(*Node_Collection); ok {
		return x.Collection
	}
	return nil
}

func (x *Node) GetEntity() *Entity {
	if x, ok := x.GetConfig().(*Node_Entity); ok {
		return x.Entity
	}
	return nil
}

func (x *Node) GetInput() *Input {
	if x, ok := x.GetConfig().(*Node_Input); ok {
		return x.Input
	}
	return nil
}

func (x *Node) GetResourceIds() []string {
	if x != nil {
		return x.ResourceIds
	}
	return nil
}

type isNode_Config interface {
	isNode_Config()
}

type Node_Rest struct {
	Rest *REST `protobuf:"bytes,6,opt,name=rest,proto3,oneof"`
}

type Node_Grpc struct {
	Grpc *GRPC `protobuf:"bytes,7,opt,name=grpc,proto3,oneof"`
}

type Node_Collection struct {
	Collection *Collection `protobuf:"bytes,8,opt,name=collection,proto3,oneof"`
}

type Node_Entity struct {
	Entity *Entity `protobuf:"bytes,9,opt,name=entity,proto3,oneof"`
}

type Node_Input struct {
	Input *Input `protobuf:"bytes,10,opt,name=input,proto3,oneof"`
}

func (*Node_Rest) isNode_Config() {}

func (*Node_Grpc) isNode_Config() {}

func (*Node_Collection) isNode_Config() {}

func (*Node_Entity) isNode_Config() {}

func (*Node_Input) isNode_Config() {}

type Edge struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	From string `protobuf:"bytes,2,opt,name=from,proto3" json:"from,omitempty"`
	To   string `protobuf:"bytes,3,opt,name=to,proto3" json:"to,omitempty"`
}

func (x *Edge) Reset() {
	*x = Edge{}
	if protoimpl.UnsafeEnabled {
		mi := &file_graph_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Edge) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Edge) ProtoMessage() {}

func (x *Edge) ProtoReflect() protoreflect.Message {
	mi := &file_graph_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Edge.ProtoReflect.Descriptor instead.
func (*Edge) Descriptor() ([]byte, []int) {
	return file_graph_proto_rawDescGZIP(), []int{2}
}

func (x *Edge) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Edge) GetFrom() string {
	if x != nil {
		return x.From
	}
	return ""
}

func (x *Edge) GetTo() string {
	if x != nil {
		return x.To
	}
	return ""
}

var File_graph_proto protoreflect.FileDescriptor

var file_graph_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x67, 0x72, 0x61, 0x70, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x67,
	0x72, 0x61, 0x70, 0x68, 0x1a, 0x0b, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x0e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x71, 0x0a, 0x05, 0x47, 0x72, 0x61, 0x70, 0x68, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x21,
	0x0a, 0x05, 0x6e, 0x6f, 0x64, 0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0b, 0x2e,
	0x67, 0x72, 0x61, 0x70, 0x68, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x05, 0x6e, 0x6f, 0x64, 0x65,
	0x73, 0x12, 0x21, 0x0a, 0x05, 0x65, 0x64, 0x67, 0x65, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x0b, 0x2e, 0x67, 0x72, 0x61, 0x70, 0x68, 0x2e, 0x45, 0x64, 0x67, 0x65, 0x52, 0x05, 0x65,
	0x64, 0x67, 0x65, 0x73, 0x22, 0xd8, 0x02, 0x0a, 0x04, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x19, 0x0a, 0x08, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x49, 0x64, 0x12, 0x0c, 0x0a, 0x01,
	0x78, 0x18, 0x04, 0x20, 0x01, 0x28, 0x02, 0x52, 0x01, 0x78, 0x12, 0x0c, 0x0a, 0x01, 0x79, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x02, 0x52, 0x01, 0x79, 0x12, 0x21, 0x0a, 0x04, 0x72, 0x65, 0x73, 0x74,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x2e, 0x52,
	0x45, 0x53, 0x54, 0x48, 0x00, 0x52, 0x04, 0x72, 0x65, 0x73, 0x74, 0x12, 0x21, 0x0a, 0x04, 0x67,
	0x72, 0x70, 0x63, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x62, 0x6c, 0x6f, 0x63,
	0x6b, 0x2e, 0x47, 0x52, 0x50, 0x43, 0x48, 0x00, 0x52, 0x04, 0x67, 0x72, 0x70, 0x63, 0x12, 0x33,
	0x0a, 0x0a, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x08, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x11, 0x2e, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x2e, 0x43, 0x6f, 0x6c, 0x6c, 0x65,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x00, 0x52, 0x0a, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x27, 0x0a, 0x06, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x18, 0x09, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x2e, 0x45, 0x6e, 0x74, 0x69,
	0x74, 0x79, 0x48, 0x00, 0x52, 0x06, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x12, 0x24, 0x0a, 0x05,
	0x69, 0x6e, 0x70, 0x75, 0x74, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x62, 0x6c,
	0x6f, 0x63, 0x6b, 0x2e, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x48, 0x00, 0x52, 0x05, 0x69, 0x6e, 0x70,
	0x75, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x69,
	0x64, 0x73, 0x18, 0x0b, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0b, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x49, 0x64, 0x73, 0x42, 0x08, 0x0a, 0x06, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x22,
	0x3a, 0x0a, 0x04, 0x45, 0x64, 0x67, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x12, 0x0e, 0x0a, 0x02, 0x74,
	0x6f, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x74, 0x6f, 0x42, 0x74, 0x0a, 0x09, 0x63,
	0x6f, 0x6d, 0x2e, 0x67, 0x72, 0x61, 0x70, 0x68, 0x42, 0x0a, 0x47, 0x72, 0x61, 0x70, 0x68, 0x50,
	0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x27, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x66, 0x6c, 0x6f, 0x77, 0x2d, 0x6c, 0x61, 0x62,
	0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x66, 0x6c, 0x6f, 0x77, 0x2f, 0x67, 0x65, 0x6e, 0xa2,
	0x02, 0x03, 0x47, 0x58, 0x58, 0xaa, 0x02, 0x05, 0x47, 0x72, 0x61, 0x70, 0x68, 0xca, 0x02, 0x05,
	0x47, 0x72, 0x61, 0x70, 0x68, 0xe2, 0x02, 0x11, 0x47, 0x72, 0x61, 0x70, 0x68, 0x5c, 0x47, 0x50,
	0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x05, 0x47, 0x72, 0x61, 0x70,
	0x68, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_graph_proto_rawDescOnce sync.Once
	file_graph_proto_rawDescData = file_graph_proto_rawDesc
)

func file_graph_proto_rawDescGZIP() []byte {
	file_graph_proto_rawDescOnce.Do(func() {
		file_graph_proto_rawDescData = protoimpl.X.CompressGZIP(file_graph_proto_rawDescData)
	})
	return file_graph_proto_rawDescData
}

var file_graph_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_graph_proto_goTypes = []interface{}{
	(*Graph)(nil),      // 0: graph.Graph
	(*Node)(nil),       // 1: graph.Node
	(*Edge)(nil),       // 2: graph.Edge
	(*REST)(nil),       // 3: block.REST
	(*GRPC)(nil),       // 4: block.GRPC
	(*Collection)(nil), // 5: block.Collection
	(*Entity)(nil),     // 6: block.Entity
	(*Input)(nil),      // 7: block.Input
}
var file_graph_proto_depIdxs = []int32{
	1, // 0: graph.Graph.nodes:type_name -> graph.Node
	2, // 1: graph.Graph.edges:type_name -> graph.Edge
	3, // 2: graph.Node.rest:type_name -> block.REST
	4, // 3: graph.Node.grpc:type_name -> block.GRPC
	5, // 4: graph.Node.collection:type_name -> block.Collection
	6, // 5: graph.Node.entity:type_name -> block.Entity
	7, // 6: graph.Node.input:type_name -> block.Input
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_graph_proto_init() }
func file_graph_proto_init() {
	if File_graph_proto != nil {
		return
	}
	file_block_proto_init()
	file_resource_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_graph_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Graph); i {
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
		file_graph_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Node); i {
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
		file_graph_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Edge); i {
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
	file_graph_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*Node_Rest)(nil),
		(*Node_Grpc)(nil),
		(*Node_Collection)(nil),
		(*Node_Entity)(nil),
		(*Node_Input)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_graph_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_graph_proto_goTypes,
		DependencyIndexes: file_graph_proto_depIdxs,
		MessageInfos:      file_graph_proto_msgTypes,
	}.Build()
	File_graph_proto = out.File
	file_graph_proto_rawDesc = nil
	file_graph_proto_goTypes = nil
	file_graph_proto_depIdxs = nil
}
