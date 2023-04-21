// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.20.3
// source: workflow.proto

package workflow

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

type ID struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *ID) Reset() {
	*x = ID{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workflow_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ID) ProtoMessage() {}

func (x *ID) ProtoReflect() protoreflect.Message {
	mi := &file_workflow_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ID.ProtoReflect.Descriptor instead.
func (*ID) Descriptor() ([]byte, []int) {
	return file_workflow_proto_rawDescGZIP(), []int{0}
}

func (x *ID) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type Workflow struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        string      `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name      string      `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Nodes     []*Node     `protobuf:"bytes,3,rep,name=nodes,proto3" json:"nodes,omitempty"`
	Edges     []*Edge     `protobuf:"bytes,4,rep,name=edges,proto3" json:"edges,omitempty"`
	Resources []*Resource `protobuf:"bytes,5,rep,name=resources,proto3" json:"resources,omitempty"`
}

func (x *Workflow) Reset() {
	*x = Workflow{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workflow_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Workflow) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Workflow) ProtoMessage() {}

func (x *Workflow) ProtoReflect() protoreflect.Message {
	mi := &file_workflow_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Workflow.ProtoReflect.Descriptor instead.
func (*Workflow) Descriptor() ([]byte, []int) {
	return file_workflow_proto_rawDescGZIP(), []int{1}
}

func (x *Workflow) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Workflow) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Workflow) GetNodes() []*Node {
	if x != nil {
		return x.Nodes
	}
	return nil
}

func (x *Workflow) GetEdges() []*Edge {
	if x != nil {
		return x.Edges
	}
	return nil
}

func (x *Workflow) GetResources() []*Resource {
	if x != nil {
		return x.Resources
	}
	return nil
}

type Resource struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Type:
	//
	//	*Resource_Db
	//	*Resource_Docstore
	//	*Resource_Bucket
	Type isResource_Type `protobuf_oneof:"type"`
}

func (x *Resource) Reset() {
	*x = Resource{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workflow_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Resource) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Resource) ProtoMessage() {}

func (x *Resource) ProtoReflect() protoreflect.Message {
	mi := &file_workflow_proto_msgTypes[2]
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
	return file_workflow_proto_rawDescGZIP(), []int{2}
}

func (m *Resource) GetType() isResource_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (x *Resource) GetDb() *DBResource {
	if x, ok := x.GetType().(*Resource_Db); ok {
		return x.Db
	}
	return nil
}

func (x *Resource) GetDocstore() *DocStoreResource {
	if x, ok := x.GetType().(*Resource_Docstore); ok {
		return x.Docstore
	}
	return nil
}

func (x *Resource) GetBucket() *BucketResource {
	if x, ok := x.GetType().(*Resource_Bucket); ok {
		return x.Bucket
	}
	return nil
}

type isResource_Type interface {
	isResource_Type()
}

type Resource_Db struct {
	Db *DBResource `protobuf:"bytes,1,opt,name=db,proto3,oneof"`
}

type Resource_Docstore struct {
	Docstore *DocStoreResource `protobuf:"bytes,2,opt,name=docstore,proto3,oneof"`
}

type Resource_Bucket struct {
	Bucket *BucketResource `protobuf:"bytes,3,opt,name=bucket,proto3,oneof"`
}

func (*Resource_Db) isResource_Type() {}

func (*Resource_Docstore) isResource_Type() {}

func (*Resource_Bucket) isResource_Type() {}

type Function struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Type:
	//
	//	*Function_Code
	Type isFunction_Type `protobuf_oneof:"type"`
}

func (x *Function) Reset() {
	*x = Function{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workflow_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Function) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Function) ProtoMessage() {}

func (x *Function) ProtoReflect() protoreflect.Message {
	mi := &file_workflow_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Function.ProtoReflect.Descriptor instead.
func (*Function) Descriptor() ([]byte, []int) {
	return file_workflow_proto_rawDescGZIP(), []int{3}
}

func (m *Function) GetType() isFunction_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (x *Function) GetCode() *CodeFunction {
	if x, ok := x.GetType().(*Function_Code); ok {
		return x.Code
	}
	return nil
}

type isFunction_Type interface {
	isFunction_Type()
}

type Function_Code struct {
	Code *CodeFunction `protobuf:"bytes,1,opt,name=code,proto3,oneof"`
}

func (*Function_Code) isFunction_Type() {}

type Input struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Params map[string]string `protobuf:"bytes,1,rep,name=params,proto3" json:"params,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Input) Reset() {
	*x = Input{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workflow_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Input) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Input) ProtoMessage() {}

func (x *Input) ProtoReflect() protoreflect.Message {
	mi := &file_workflow_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Input.ProtoReflect.Descriptor instead.
func (*Input) Descriptor() ([]byte, []int) {
	return file_workflow_proto_rawDescGZIP(), []int{4}
}

func (x *Input) GetParams() map[string]string {
	if x != nil {
		return x.Params
	}
	return nil
}

type Collection struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *Collection) Reset() {
	*x = Collection{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workflow_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Collection) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Collection) ProtoMessage() {}

func (x *Collection) ProtoReflect() protoreflect.Message {
	mi := &file_workflow_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Collection.ProtoReflect.Descriptor instead.
func (*Collection) Descriptor() ([]byte, []int) {
	return file_workflow_proto_rawDescGZIP(), []int{5}
}

func (x *Collection) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type Table struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *Table) Reset() {
	*x = Table{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workflow_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Table) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Table) ProtoMessage() {}

func (x *Table) ProtoReflect() protoreflect.Message {
	mi := &file_workflow_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Table.ProtoReflect.Descriptor instead.
func (*Table) Descriptor() ([]byte, []int) {
	return file_workflow_proto_rawDescGZIP(), []int{6}
}

func (x *Table) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type Data struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Type:
	//
	//	*Data_Input
	//	*Data_Collection
	Type isData_Type `protobuf_oneof:"type"`
}

func (x *Data) Reset() {
	*x = Data{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workflow_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Data) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Data) ProtoMessage() {}

func (x *Data) ProtoReflect() protoreflect.Message {
	mi := &file_workflow_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Data.ProtoReflect.Descriptor instead.
func (*Data) Descriptor() ([]byte, []int) {
	return file_workflow_proto_rawDescGZIP(), []int{7}
}

func (m *Data) GetType() isData_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (x *Data) GetInput() *Input {
	if x, ok := x.GetType().(*Data_Input); ok {
		return x.Input
	}
	return nil
}

func (x *Data) GetCollection() *Collection {
	if x, ok := x.GetType().(*Data_Collection); ok {
		return x.Collection
	}
	return nil
}

type isData_Type interface {
	isData_Type()
}

type Data_Input struct {
	Input *Input `protobuf:"bytes,1,opt,name=input,proto3,oneof"`
}

type Data_Collection struct {
	Collection *Collection `protobuf:"bytes,2,opt,name=collection,proto3,oneof"`
}

func (*Data_Input) isData_Type() {}

func (*Data_Collection) isData_Type() {}

type Node struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// Types that are assignable to Type:
	//
	//	*Node_Function
	//	*Node_Data
	Type isNode_Type `protobuf_oneof:"type"`
}

func (x *Node) Reset() {
	*x = Node{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workflow_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Node) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Node) ProtoMessage() {}

func (x *Node) ProtoReflect() protoreflect.Message {
	mi := &file_workflow_proto_msgTypes[8]
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
	return file_workflow_proto_rawDescGZIP(), []int{8}
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

func (m *Node) GetType() isNode_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (x *Node) GetFunction() *Function {
	if x, ok := x.GetType().(*Node_Function); ok {
		return x.Function
	}
	return nil
}

func (x *Node) GetData() *Data {
	if x, ok := x.GetType().(*Node_Data); ok {
		return x.Data
	}
	return nil
}

type isNode_Type interface {
	isNode_Type()
}

type Node_Function struct {
	Function *Function `protobuf:"bytes,3,opt,name=function,proto3,oneof"`
}

type Node_Data struct {
	Data *Data `protobuf:"bytes,4,opt,name=data,proto3,oneof"`
}

func (*Node_Function) isNode_Type() {}

func (*Node_Data) isNode_Type() {}

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
		mi := &file_workflow_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Edge) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Edge) ProtoMessage() {}

func (x *Edge) ProtoReflect() protoreflect.Message {
	mi := &file_workflow_proto_msgTypes[9]
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
	return file_workflow_proto_rawDescGZIP(), []int{9}
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

type GRPCFunction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Host    string `protobuf:"bytes,1,opt,name=host,proto3" json:"host,omitempty"`
	Service string `protobuf:"bytes,2,opt,name=service,proto3" json:"service,omitempty"`
	Method  string `protobuf:"bytes,3,opt,name=method,proto3" json:"method,omitempty"`
}

func (x *GRPCFunction) Reset() {
	*x = GRPCFunction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workflow_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GRPCFunction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GRPCFunction) ProtoMessage() {}

func (x *GRPCFunction) ProtoReflect() protoreflect.Message {
	mi := &file_workflow_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GRPCFunction.ProtoReflect.Descriptor instead.
func (*GRPCFunction) Descriptor() ([]byte, []int) {
	return file_workflow_proto_rawDescGZIP(), []int{10}
}

func (x *GRPCFunction) GetHost() string {
	if x != nil {
		return x.Host
	}
	return ""
}

func (x *GRPCFunction) GetService() string {
	if x != nil {
		return x.Service
	}
	return ""
}

func (x *GRPCFunction) GetMethod() string {
	if x != nil {
		return x.Method
	}
	return ""
}

type CodeFunction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code string `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
}

func (x *CodeFunction) Reset() {
	*x = CodeFunction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workflow_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CodeFunction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CodeFunction) ProtoMessage() {}

func (x *CodeFunction) ProtoReflect() protoreflect.Message {
	mi := &file_workflow_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CodeFunction.ProtoReflect.Descriptor instead.
func (*CodeFunction) Descriptor() ([]byte, []int) {
	return file_workflow_proto_rawDescGZIP(), []int{11}
}

func (x *CodeFunction) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

type SQLFunction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url   string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	Query string `protobuf:"bytes,2,opt,name=query,proto3" json:"query,omitempty"`
}

func (x *SQLFunction) Reset() {
	*x = SQLFunction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workflow_proto_msgTypes[12]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SQLFunction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SQLFunction) ProtoMessage() {}

func (x *SQLFunction) ProtoReflect() protoreflect.Message {
	mi := &file_workflow_proto_msgTypes[12]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SQLFunction.ProtoReflect.Descriptor instead.
func (*SQLFunction) Descriptor() ([]byte, []int) {
	return file_workflow_proto_rawDescGZIP(), []int{12}
}

func (x *SQLFunction) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *SQLFunction) GetQuery() string {
	if x != nil {
		return x.Query
	}
	return ""
}

type DBResource struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
}

func (x *DBResource) Reset() {
	*x = DBResource{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workflow_proto_msgTypes[13]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DBResource) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DBResource) ProtoMessage() {}

func (x *DBResource) ProtoReflect() protoreflect.Message {
	mi := &file_workflow_proto_msgTypes[13]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DBResource.ProtoReflect.Descriptor instead.
func (*DBResource) Descriptor() ([]byte, []int) {
	return file_workflow_proto_rawDescGZIP(), []int{13}
}

func (x *DBResource) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

type DocStoreResource struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
}

func (x *DocStoreResource) Reset() {
	*x = DocStoreResource{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workflow_proto_msgTypes[14]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DocStoreResource) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DocStoreResource) ProtoMessage() {}

func (x *DocStoreResource) ProtoReflect() protoreflect.Message {
	mi := &file_workflow_proto_msgTypes[14]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DocStoreResource.ProtoReflect.Descriptor instead.
func (*DocStoreResource) Descriptor() ([]byte, []int) {
	return file_workflow_proto_rawDescGZIP(), []int{14}
}

func (x *DocStoreResource) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

type BucketResource struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
}

func (x *BucketResource) Reset() {
	*x = BucketResource{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workflow_proto_msgTypes[15]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BucketResource) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BucketResource) ProtoMessage() {}

func (x *BucketResource) ProtoReflect() protoreflect.Message {
	mi := &file_workflow_proto_msgTypes[15]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BucketResource.ProtoReflect.Descriptor instead.
func (*BucketResource) Descriptor() ([]byte, []int) {
	return file_workflow_proto_rawDescGZIP(), []int{15}
}

func (x *BucketResource) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

var File_workflow_proto protoreflect.FileDescriptor

var file_workflow_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x08, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x22, 0x14, 0x0a, 0x02, 0x49, 0x44,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64,
	0x22, 0xac, 0x01, 0x0a, 0x08, 0x57, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x24, 0x0a, 0x05, 0x6e, 0x6f, 0x64, 0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x0e, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x2e, 0x4e, 0x6f, 0x64, 0x65,
	0x52, 0x05, 0x6e, 0x6f, 0x64, 0x65, 0x73, 0x12, 0x24, 0x0a, 0x05, 0x65, 0x64, 0x67, 0x65, 0x73,
	0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f,
	0x77, 0x2e, 0x45, 0x64, 0x67, 0x65, 0x52, 0x05, 0x65, 0x64, 0x67, 0x65, 0x73, 0x12, 0x30, 0x0a,
	0x09, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x12, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x2e, 0x52, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x52, 0x09, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x22,
	0xa8, 0x01, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x26, 0x0a, 0x02,
	0x64, 0x62, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x66,
	0x6c, 0x6f, 0x77, 0x2e, 0x44, 0x42, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x48, 0x00,
	0x52, 0x02, 0x64, 0x62, 0x12, 0x38, 0x0a, 0x08, 0x64, 0x6f, 0x63, 0x73, 0x74, 0x6f, 0x72, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f,
	0x77, 0x2e, 0x44, 0x6f, 0x63, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x48, 0x00, 0x52, 0x08, 0x64, 0x6f, 0x63, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x12, 0x32,
	0x0a, 0x06, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18,
	0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x2e, 0x42, 0x75, 0x63, 0x6b, 0x65, 0x74,
	0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x48, 0x00, 0x52, 0x06, 0x62, 0x75, 0x63, 0x6b,
	0x65, 0x74, 0x42, 0x06, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x40, 0x0a, 0x08, 0x46, 0x75,
	0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x2c, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x2e,
	0x43, 0x6f, 0x64, 0x65, 0x46, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x00, 0x52, 0x04,
	0x63, 0x6f, 0x64, 0x65, 0x42, 0x06, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x77, 0x0a, 0x05,
	0x49, 0x6e, 0x70, 0x75, 0x74, 0x12, 0x33, 0x0a, 0x06, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77,
	0x2e, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x2e, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x52, 0x06, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x1a, 0x39, 0x0a, 0x0b, 0x50, 0x61,
	0x72, 0x61, 0x6d, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x20, 0x0a, 0x0a, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x1b, 0x0a, 0x05, 0x54, 0x61, 0x62, 0x6c, 0x65,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x22, 0x6f, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x12, 0x27, 0x0a, 0x05,
	0x69, 0x6e, 0x70, 0x75, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x77, 0x6f,
	0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x2e, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x48, 0x00, 0x52, 0x05,
	0x69, 0x6e, 0x70, 0x75, 0x74, 0x12, 0x36, 0x0a, 0x0a, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x77, 0x6f, 0x72, 0x6b,
	0x66, 0x6c, 0x6f, 0x77, 0x2e, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x48,
	0x00, 0x52, 0x0a, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x06, 0x0a,
	0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x8a, 0x01, 0x0a, 0x04, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12,
	0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x30, 0x0a, 0x08, 0x66, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x2e,
	0x46, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x00, 0x52, 0x08, 0x66, 0x75, 0x6e, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x24, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x2e, 0x44, 0x61,
	0x74, 0x61, 0x48, 0x00, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x42, 0x06, 0x0a, 0x04, 0x74, 0x79,
	0x70, 0x65, 0x22, 0x3a, 0x0a, 0x04, 0x45, 0x64, 0x67, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x66, 0x72,
	0x6f, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x12, 0x0e,
	0x0a, 0x02, 0x74, 0x6f, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x74, 0x6f, 0x22, 0x54,
	0x0a, 0x0c, 0x47, 0x52, 0x50, 0x43, 0x46, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x12,
	0x0a, 0x04, 0x68, 0x6f, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x6f,
	0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x16, 0x0a, 0x06,
	0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6d, 0x65,
	0x74, 0x68, 0x6f, 0x64, 0x22, 0x22, 0x0a, 0x0c, 0x43, 0x6f, 0x64, 0x65, 0x46, 0x75, 0x6e, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x22, 0x35, 0x0a, 0x0b, 0x53, 0x51, 0x4c, 0x46,
	0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x12, 0x14, 0x0a, 0x05, 0x71, 0x75, 0x65,
	0x72, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x71, 0x75, 0x65, 0x72, 0x79, 0x22,
	0x1e, 0x0a, 0x0a, 0x44, 0x42, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x10, 0x0a,
	0x03, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x22,
	0x24, 0x0a, 0x10, 0x44, 0x6f, 0x63, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x52, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x75, 0x72, 0x6c, 0x22, 0x22, 0x0a, 0x0e, 0x42, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x52,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x32, 0x3f, 0x0a, 0x07, 0x4d, 0x61, 0x6e,
	0x61, 0x67, 0x65, 0x72, 0x12, 0x34, 0x0a, 0x0e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x57, 0x6f,
	0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x12, 0x12, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f,
	0x77, 0x2e, 0x57, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x1a, 0x0c, 0x2e, 0x77, 0x6f, 0x72,
	0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x2e, 0x49, 0x44, 0x22, 0x00, 0x42, 0x10, 0x5a, 0x0e, 0x2e, 0x2f,
	0x67, 0x65, 0x6e, 0x2f, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_workflow_proto_rawDescOnce sync.Once
	file_workflow_proto_rawDescData = file_workflow_proto_rawDesc
)

func file_workflow_proto_rawDescGZIP() []byte {
	file_workflow_proto_rawDescOnce.Do(func() {
		file_workflow_proto_rawDescData = protoimpl.X.CompressGZIP(file_workflow_proto_rawDescData)
	})
	return file_workflow_proto_rawDescData
}

var file_workflow_proto_msgTypes = make([]protoimpl.MessageInfo, 17)
var file_workflow_proto_goTypes = []interface{}{
	(*ID)(nil),               // 0: workflow.ID
	(*Workflow)(nil),         // 1: workflow.Workflow
	(*Resource)(nil),         // 2: workflow.Resource
	(*Function)(nil),         // 3: workflow.Function
	(*Input)(nil),            // 4: workflow.Input
	(*Collection)(nil),       // 5: workflow.Collection
	(*Table)(nil),            // 6: workflow.Table
	(*Data)(nil),             // 7: workflow.Data
	(*Node)(nil),             // 8: workflow.Node
	(*Edge)(nil),             // 9: workflow.Edge
	(*GRPCFunction)(nil),     // 10: workflow.GRPCFunction
	(*CodeFunction)(nil),     // 11: workflow.CodeFunction
	(*SQLFunction)(nil),      // 12: workflow.SQLFunction
	(*DBResource)(nil),       // 13: workflow.DBResource
	(*DocStoreResource)(nil), // 14: workflow.DocStoreResource
	(*BucketResource)(nil),   // 15: workflow.BucketResource
	nil,                      // 16: workflow.Input.ParamsEntry
}
var file_workflow_proto_depIdxs = []int32{
	8,  // 0: workflow.Workflow.nodes:type_name -> workflow.Node
	9,  // 1: workflow.Workflow.edges:type_name -> workflow.Edge
	2,  // 2: workflow.Workflow.resources:type_name -> workflow.Resource
	13, // 3: workflow.Resource.db:type_name -> workflow.DBResource
	14, // 4: workflow.Resource.docstore:type_name -> workflow.DocStoreResource
	15, // 5: workflow.Resource.bucket:type_name -> workflow.BucketResource
	11, // 6: workflow.Function.code:type_name -> workflow.CodeFunction
	16, // 7: workflow.Input.params:type_name -> workflow.Input.ParamsEntry
	4,  // 8: workflow.Data.input:type_name -> workflow.Input
	5,  // 9: workflow.Data.collection:type_name -> workflow.Collection
	3,  // 10: workflow.Node.function:type_name -> workflow.Function
	7,  // 11: workflow.Node.data:type_name -> workflow.Data
	1,  // 12: workflow.Manager.CreateWorkflow:input_type -> workflow.Workflow
	0,  // 13: workflow.Manager.CreateWorkflow:output_type -> workflow.ID
	13, // [13:14] is the sub-list for method output_type
	12, // [12:13] is the sub-list for method input_type
	12, // [12:12] is the sub-list for extension type_name
	12, // [12:12] is the sub-list for extension extendee
	0,  // [0:12] is the sub-list for field type_name
}

func init() { file_workflow_proto_init() }
func file_workflow_proto_init() {
	if File_workflow_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_workflow_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ID); i {
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
		file_workflow_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Workflow); i {
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
		file_workflow_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
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
		file_workflow_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Function); i {
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
		file_workflow_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Input); i {
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
		file_workflow_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Collection); i {
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
		file_workflow_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Table); i {
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
		file_workflow_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Data); i {
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
		file_workflow_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
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
		file_workflow_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
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
		file_workflow_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GRPCFunction); i {
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
		file_workflow_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CodeFunction); i {
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
		file_workflow_proto_msgTypes[12].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SQLFunction); i {
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
		file_workflow_proto_msgTypes[13].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DBResource); i {
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
		file_workflow_proto_msgTypes[14].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DocStoreResource); i {
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
		file_workflow_proto_msgTypes[15].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BucketResource); i {
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
	file_workflow_proto_msgTypes[2].OneofWrappers = []interface{}{
		(*Resource_Db)(nil),
		(*Resource_Docstore)(nil),
		(*Resource_Bucket)(nil),
	}
	file_workflow_proto_msgTypes[3].OneofWrappers = []interface{}{
		(*Function_Code)(nil),
	}
	file_workflow_proto_msgTypes[7].OneofWrappers = []interface{}{
		(*Data_Input)(nil),
		(*Data_Collection)(nil),
	}
	file_workflow_proto_msgTypes[8].OneofWrappers = []interface{}{
		(*Node_Function)(nil),
		(*Node_Data)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_workflow_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   17,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_workflow_proto_goTypes,
		DependencyIndexes: file_workflow_proto_depIdxs,
		MessageInfos:      file_workflow_proto_msgTypes,
	}.Build()
	File_workflow_proto = out.File
	file_workflow_proto_rawDesc = nil
	file_workflow_proto_goTypes = nil
	file_workflow_proto_depIdxs = nil
}
