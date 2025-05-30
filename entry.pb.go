// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v3.21.12
// source: entry.proto

package walle

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

type Entry struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Lsn           int64                  `protobuf:"varint,1,opt,name=lsn,proto3" json:"lsn,omitempty"`
	Data          []byte                 `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	Checksum      uint32                 `protobuf:"varint,3,opt,name=checksum,proto3" json:"checksum,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Entry) Reset() {
	*x = Entry{}
	mi := &file_entry_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Entry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Entry) ProtoMessage() {}

func (x *Entry) ProtoReflect() protoreflect.Message {
	mi := &file_entry_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Entry.ProtoReflect.Descriptor instead.
func (*Entry) Descriptor() ([]byte, []int) {
	return file_entry_proto_rawDescGZIP(), []int{0}
}

func (x *Entry) GetLsn() int64 {
	if x != nil {
		return x.Lsn
	}
	return 0
}

func (x *Entry) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *Entry) GetChecksum() uint32 {
	if x != nil {
		return x.Checksum
	}
	return 0
}

var File_entry_proto protoreflect.FileDescriptor

const file_entry_proto_rawDesc = "" +
	"\n" +
	"\ventry.proto\x12\x05walle\"I\n" +
	"\x05Entry\x12\x10\n" +
	"\x03lsn\x18\x01 \x01(\x03R\x03lsn\x12\x12\n" +
	"\x04data\x18\x02 \x01(\fR\x04data\x12\x1a\n" +
	"\bchecksum\x18\x03 \x01(\rR\bchecksumB\x03Z\x01.b\x06proto3"

var (
	file_entry_proto_rawDescOnce sync.Once
	file_entry_proto_rawDescData []byte
)

func file_entry_proto_rawDescGZIP() []byte {
	file_entry_proto_rawDescOnce.Do(func() {
		file_entry_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_entry_proto_rawDesc), len(file_entry_proto_rawDesc)))
	})
	return file_entry_proto_rawDescData
}

var file_entry_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_entry_proto_goTypes = []any{
	(*Entry)(nil), // 0: walle.Entry
}
var file_entry_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_entry_proto_init() }
func file_entry_proto_init() {
	if File_entry_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_entry_proto_rawDesc), len(file_entry_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_entry_proto_goTypes,
		DependencyIndexes: file_entry_proto_depIdxs,
		MessageInfos:      file_entry_proto_msgTypes,
	}.Build()
	File_entry_proto = out.File
	file_entry_proto_goTypes = nil
	file_entry_proto_depIdxs = nil
}
