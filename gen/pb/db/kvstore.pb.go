// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v3.21.12
// source: db/kvstore.proto

package db

import (
	graph "github.com/gemlab-dev/relor/gen/pb/graph"
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

type WorkflowStatus int32

const (
	WorkflowStatus_WORKFLOW_STATUS_UNSPECIFIED WorkflowStatus = 0
	WorkflowStatus_WORKFLOW_STATUS_PENDING     WorkflowStatus = 1
	WorkflowStatus_WORKFLOW_STATUS_RUNNING     WorkflowStatus = 2
	WorkflowStatus_WORKFLOW_STATUS_COMPLETED   WorkflowStatus = 3
	WorkflowStatus_WORKFLOW_STATUS_FAILED      WorkflowStatus = 4
)

// Enum value maps for WorkflowStatus.
var (
	WorkflowStatus_name = map[int32]string{
		0: "WORKFLOW_STATUS_UNSPECIFIED",
		1: "WORKFLOW_STATUS_PENDING",
		2: "WORKFLOW_STATUS_RUNNING",
		3: "WORKFLOW_STATUS_COMPLETED",
		4: "WORKFLOW_STATUS_FAILED",
	}
	WorkflowStatus_value = map[string]int32{
		"WORKFLOW_STATUS_UNSPECIFIED": 0,
		"WORKFLOW_STATUS_PENDING":     1,
		"WORKFLOW_STATUS_RUNNING":     2,
		"WORKFLOW_STATUS_COMPLETED":   3,
		"WORKFLOW_STATUS_FAILED":      4,
	}
)

func (x WorkflowStatus) Enum() *WorkflowStatus {
	p := new(WorkflowStatus)
	*p = x
	return p
}

func (x WorkflowStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (WorkflowStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_db_kvstore_proto_enumTypes[0].Descriptor()
}

func (WorkflowStatus) Type() protoreflect.EnumType {
	return &file_db_kvstore_proto_enumTypes[0]
}

func (x WorkflowStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use WorkflowStatus.Descriptor instead.
func (WorkflowStatus) EnumDescriptor() ([]byte, []int) {
	return file_db_kvstore_proto_rawDescGZIP(), []int{0}
}

type Workflow struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          string         `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	CurrentNode string         `protobuf:"bytes,2,opt,name=current_node,json=currentNode,proto3" json:"current_node,omitempty"`
	Status      WorkflowStatus `protobuf:"varint,3,opt,name=status,proto3,enum=kvstore.WorkflowStatus" json:"status,omitempty"`
	Graph       *graph.Graph   `protobuf:"bytes,4,opt,name=graph,proto3" json:"graph,omitempty"`
}

func (x *Workflow) Reset() {
	*x = Workflow{}
	mi := &file_db_kvstore_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Workflow) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Workflow) ProtoMessage() {}

func (x *Workflow) ProtoReflect() protoreflect.Message {
	mi := &file_db_kvstore_proto_msgTypes[0]
	if x != nil {
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
	return file_db_kvstore_proto_rawDescGZIP(), []int{0}
}

func (x *Workflow) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Workflow) GetCurrentNode() string {
	if x != nil {
		return x.CurrentNode
	}
	return ""
}

func (x *Workflow) GetStatus() WorkflowStatus {
	if x != nil {
		return x.Status
	}
	return WorkflowStatus_WORKFLOW_STATUS_UNSPECIFIED
}

func (x *Workflow) GetGraph() *graph.Graph {
	if x != nil {
		return x.Graph
	}
	return nil
}

var File_db_kvstore_proto protoreflect.FileDescriptor

var file_db_kvstore_proto_rawDesc = []byte{
	0x0a, 0x10, 0x64, 0x62, 0x2f, 0x6b, 0x76, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x07, 0x6b, 0x76, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x1a, 0x11, 0x67, 0x72, 0x61,
	0x70, 0x68, 0x2f, 0x67, 0x72, 0x61, 0x70, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x92,
	0x01, 0x0a, 0x08, 0x57, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x21, 0x0a, 0x0c, 0x63,
	0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x6e, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x2f,
	0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x17,
	0x2e, 0x6b, 0x76, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x57, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f,
	0x77, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12,
	0x22, 0x0a, 0x05, 0x67, 0x72, 0x61, 0x70, 0x68, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c,
	0x2e, 0x67, 0x72, 0x61, 0x70, 0x68, 0x2e, 0x47, 0x72, 0x61, 0x70, 0x68, 0x52, 0x05, 0x67, 0x72,
	0x61, 0x70, 0x68, 0x2a, 0xa6, 0x01, 0x0a, 0x0e, 0x57, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1f, 0x0a, 0x1b, 0x57, 0x4f, 0x52, 0x4b, 0x46, 0x4c,
	0x4f, 0x57, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43,
	0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x1b, 0x0a, 0x17, 0x57, 0x4f, 0x52, 0x4b, 0x46,
	0x4c, 0x4f, 0x57, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x50, 0x45, 0x4e, 0x44, 0x49,
	0x4e, 0x47, 0x10, 0x01, 0x12, 0x1b, 0x0a, 0x17, 0x57, 0x4f, 0x52, 0x4b, 0x46, 0x4c, 0x4f, 0x57,
	0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x52, 0x55, 0x4e, 0x4e, 0x49, 0x4e, 0x47, 0x10,
	0x02, 0x12, 0x1d, 0x0a, 0x19, 0x57, 0x4f, 0x52, 0x4b, 0x46, 0x4c, 0x4f, 0x57, 0x5f, 0x53, 0x54,
	0x41, 0x54, 0x55, 0x53, 0x5f, 0x43, 0x4f, 0x4d, 0x50, 0x4c, 0x45, 0x54, 0x45, 0x44, 0x10, 0x03,
	0x12, 0x1a, 0x0a, 0x16, 0x57, 0x4f, 0x52, 0x4b, 0x46, 0x4c, 0x4f, 0x57, 0x5f, 0x53, 0x54, 0x41,
	0x54, 0x55, 0x53, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x04, 0x42, 0x30, 0x5a, 0x2e,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x65, 0x6d, 0x6c, 0x61,
	0x62, 0x2d, 0x64, 0x65, 0x76, 0x2f, 0x72, 0x65, 0x6c, 0x6f, 0x72, 0x2f, 0x69, 0x6e, 0x74, 0x65,
	0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x62, 0x2f, 0x64, 0x62, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_db_kvstore_proto_rawDescOnce sync.Once
	file_db_kvstore_proto_rawDescData = file_db_kvstore_proto_rawDesc
)

func file_db_kvstore_proto_rawDescGZIP() []byte {
	file_db_kvstore_proto_rawDescOnce.Do(func() {
		file_db_kvstore_proto_rawDescData = protoimpl.X.CompressGZIP(file_db_kvstore_proto_rawDescData)
	})
	return file_db_kvstore_proto_rawDescData
}

var file_db_kvstore_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_db_kvstore_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_db_kvstore_proto_goTypes = []any{
	(WorkflowStatus)(0), // 0: kvstore.WorkflowStatus
	(*Workflow)(nil),    // 1: kvstore.Workflow
	(*graph.Graph)(nil), // 2: graph.Graph
}
var file_db_kvstore_proto_depIdxs = []int32{
	0, // 0: kvstore.Workflow.status:type_name -> kvstore.WorkflowStatus
	2, // 1: kvstore.Workflow.graph:type_name -> graph.Graph
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_db_kvstore_proto_init() }
func file_db_kvstore_proto_init() {
	if File_db_kvstore_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_db_kvstore_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_db_kvstore_proto_goTypes,
		DependencyIndexes: file_db_kvstore_proto_depIdxs,
		EnumInfos:         file_db_kvstore_proto_enumTypes,
		MessageInfos:      file_db_kvstore_proto_msgTypes,
	}.Build()
	File_db_kvstore_proto = out.File
	file_db_kvstore_proto_rawDesc = nil
	file_db_kvstore_proto_goTypes = nil
	file_db_kvstore_proto_depIdxs = nil
}
