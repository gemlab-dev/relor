// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v3.21.12
// source: config/config.proto

package config

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Generic address message.
type Address struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Hostname or IP address.
	Hostname string `protobuf:"bytes,1,opt,name=hostname,proto3" json:"hostname,omitempty"`
	// Port number.
	Port          int32 `protobuf:"varint,2,opt,name=port,proto3" json:"port,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Address) Reset() {
	*x = Address{}
	mi := &file_config_config_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Address) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Address) ProtoMessage() {}

func (x *Address) ProtoReflect() protoreflect.Message {
	mi := &file_config_config_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Address.ProtoReflect.Descriptor instead.
func (*Address) Descriptor() ([]byte, []int) {
	return file_config_config_proto_rawDescGZIP(), []int{0}
}

func (x *Address) GetHostname() string {
	if x != nil {
		return x.Hostname
	}
	return ""
}

func (x *Address) GetPort() int32 {
	if x != nil {
		return x.Port
	}
	return 0
}

// Node-to-node communication settings.
type Cluster struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Unique identifier of the instace of the service.
	NodeName string `protobuf:"bytes,1,opt,name=node_name,json=nodeName,proto3" json:"node_name,omitempty"`
	// Port the gossip protocol is exposed on.
	GossipPort int32 `protobuf:"varint,2,opt,name=gossip_port,json=gossipPort,proto3" json:"gossip_port,omitempty"`
	// List of seed nodes to connect to.
	GossipSeed    []*Address `protobuf:"bytes,3,rep,name=gossip_seed,json=gossipSeed,proto3" json:"gossip_seed,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Cluster) Reset() {
	*x = Cluster{}
	mi := &file_config_config_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Cluster) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Cluster) ProtoMessage() {}

func (x *Cluster) ProtoReflect() protoreflect.Message {
	mi := &file_config_config_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Cluster.ProtoReflect.Descriptor instead.
func (*Cluster) Descriptor() ([]byte, []int) {
	return file_config_config_proto_rawDescGZIP(), []int{1}
}

func (x *Cluster) GetNodeName() string {
	if x != nil {
		return x.NodeName
	}
	return ""
}

func (x *Cluster) GetGossipPort() int32 {
	if x != nil {
		return x.GossipPort
	}
	return 0
}

func (x *Cluster) GetGossipSeed() []*Address {
	if x != nil {
		return x.GossipSeed
	}
	return nil
}

type Storage struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Path          string                 `protobuf:"bytes,1,opt,name=path,proto3" json:"path,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Storage) Reset() {
	*x = Storage{}
	mi := &file_config_config_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Storage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Storage) ProtoMessage() {}

func (x *Storage) ProtoReflect() protoreflect.Message {
	mi := &file_config_config_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Storage.ProtoReflect.Descriptor instead.
func (*Storage) Descriptor() ([]byte, []int) {
	return file_config_config_proto_rawDescGZIP(), []int{2}
}

func (x *Storage) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

// Orchestrator service configuration.
type Config struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Managament API port.
	ApiPort int32 `protobuf:"varint,1,opt,name=api_port,json=apiPort,proto3" json:"api_port,omitempty"`
	// Address exposing Jobs API.
	JobServiceAddr *Address `protobuf:"bytes,2,opt,name=job_service_addr,json=jobServiceAddr,proto3" json:"job_service_addr,omitempty"`
	// Specify settings for a cluster mode.
	// If this value is not set, single node configuration is used.
	Cluster *Cluster `protobuf:"bytes,3,opt,name=cluster,proto3,oneof" json:"cluster,omitempty"`
	// Specify settings for a storage.
	Storage       *Storage `protobuf:"bytes,4,opt,name=storage,proto3" json:"storage,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Config) Reset() {
	*x = Config{}
	mi := &file_config_config_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Config) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Config) ProtoMessage() {}

func (x *Config) ProtoReflect() protoreflect.Message {
	mi := &file_config_config_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Config.ProtoReflect.Descriptor instead.
func (*Config) Descriptor() ([]byte, []int) {
	return file_config_config_proto_rawDescGZIP(), []int{3}
}

func (x *Config) GetApiPort() int32 {
	if x != nil {
		return x.ApiPort
	}
	return 0
}

func (x *Config) GetJobServiceAddr() *Address {
	if x != nil {
		return x.JobServiceAddr
	}
	return nil
}

func (x *Config) GetCluster() *Cluster {
	if x != nil {
		return x.Cluster
	}
	return nil
}

func (x *Config) GetStorage() *Storage {
	if x != nil {
		return x.Storage
	}
	return nil
}

var File_config_config_proto protoreflect.FileDescriptor

const file_config_config_proto_rawDesc = "" +
	"\n" +
	"\x13config/config.proto\x12\x06config\"9\n" +
	"\aAddress\x12\x1a\n" +
	"\bhostname\x18\x01 \x01(\tR\bhostname\x12\x12\n" +
	"\x04port\x18\x02 \x01(\x05R\x04port\"y\n" +
	"\aCluster\x12\x1b\n" +
	"\tnode_name\x18\x01 \x01(\tR\bnodeName\x12\x1f\n" +
	"\vgossip_port\x18\x02 \x01(\x05R\n" +
	"gossipPort\x120\n" +
	"\vgossip_seed\x18\x03 \x03(\v2\x0f.config.AddressR\n" +
	"gossipSeed\"\x1d\n" +
	"\aStorage\x12\x12\n" +
	"\x04path\x18\x01 \x01(\tR\x04path\"\xc5\x01\n" +
	"\x06Config\x12\x19\n" +
	"\bapi_port\x18\x01 \x01(\x05R\aapiPort\x129\n" +
	"\x10job_service_addr\x18\x02 \x01(\v2\x0f.config.AddressR\x0ejobServiceAddr\x12.\n" +
	"\acluster\x18\x03 \x01(\v2\x0f.config.ClusterH\x00R\acluster\x88\x01\x01\x12)\n" +
	"\astorage\x18\x04 \x01(\v2\x0f.config.StorageR\astorageB\n" +
	"\n" +
	"\b_clusterB+Z)github.com/gemlab-dev/relor/gen/pb/configb\x06proto3"

var (
	file_config_config_proto_rawDescOnce sync.Once
	file_config_config_proto_rawDescData []byte
)

func file_config_config_proto_rawDescGZIP() []byte {
	file_config_config_proto_rawDescOnce.Do(func() {
		file_config_config_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_config_config_proto_rawDesc), len(file_config_config_proto_rawDesc)))
	})
	return file_config_config_proto_rawDescData
}

var file_config_config_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_config_config_proto_goTypes = []any{
	(*Address)(nil), // 0: config.Address
	(*Cluster)(nil), // 1: config.Cluster
	(*Storage)(nil), // 2: config.Storage
	(*Config)(nil),  // 3: config.Config
}
var file_config_config_proto_depIdxs = []int32{
	0, // 0: config.Cluster.gossip_seed:type_name -> config.Address
	0, // 1: config.Config.job_service_addr:type_name -> config.Address
	1, // 2: config.Config.cluster:type_name -> config.Cluster
	2, // 3: config.Config.storage:type_name -> config.Storage
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_config_config_proto_init() }
func file_config_config_proto_init() {
	if File_config_config_proto != nil {
		return
	}
	file_config_config_proto_msgTypes[3].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_config_config_proto_rawDesc), len(file_config_config_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_config_config_proto_goTypes,
		DependencyIndexes: file_config_config_proto_depIdxs,
		MessageInfos:      file_config_config_proto_msgTypes,
	}.Build()
	File_config_config_proto = out.File
	file_config_config_proto_goTypes = nil
	file_config_config_proto_depIdxs = nil
}
