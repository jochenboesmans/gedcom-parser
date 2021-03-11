// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.15.5
// source: grpc/parse.proto

package grpc

import (
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type PathsToFiles struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	InputFilePath  string `protobuf:"bytes,1,opt,name=inputFilePath,proto3" json:"inputFilePath,omitempty"`
	OutputFilePath string `protobuf:"bytes,2,opt,name=outputFilePath,proto3" json:"outputFilePath,omitempty"`
}

func (x *PathsToFiles) Reset() {
	*x = PathsToFiles{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_parse_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PathsToFiles) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PathsToFiles) ProtoMessage() {}

func (x *PathsToFiles) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_parse_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PathsToFiles.ProtoReflect.Descriptor instead.
func (*PathsToFiles) Descriptor() ([]byte, []int) {
	return file_grpc_parse_proto_rawDescGZIP(), []int{0}
}

func (x *PathsToFiles) GetInputFilePath() string {
	if x != nil {
		return x.InputFilePath
	}
	return ""
}

func (x *PathsToFiles) GetOutputFilePath() string {
	if x != nil {
		return x.OutputFilePath
	}
	return ""
}

type Result struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message   string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	Error     string `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
	Something string `protobuf:"bytes,3,opt,name=something,proto3" json:"something,omitempty"`
}

func (x *Result) Reset() {
	*x = Result{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_parse_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Result) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Result) ProtoMessage() {}

func (x *Result) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_parse_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Result.ProtoReflect.Descriptor instead.
func (*Result) Descriptor() ([]byte, []int) {
	return file_grpc_parse_proto_rawDescGZIP(), []int{1}
}

func (x *Result) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *Result) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

func (x *Result) GetSomething() string {
	if x != nil {
		return x.Something
	}
	return ""
}

var File_grpc_parse_proto protoreflect.FileDescriptor

var file_grpc_parse_proto_rawDesc = []byte{
	0x0a, 0x10, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x70, 0x61, 0x72, 0x73, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x04, 0x67, 0x72, 0x70, 0x63, 0x22, 0x5c, 0x0a, 0x0c, 0x50, 0x61, 0x74, 0x68,
	0x73, 0x54, 0x6f, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x12, 0x24, 0x0a, 0x0d, 0x69, 0x6e, 0x70, 0x75,
	0x74, 0x46, 0x69, 0x6c, 0x65, 0x50, 0x61, 0x74, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0d, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x50, 0x61, 0x74, 0x68, 0x12, 0x26,
	0x0a, 0x0e, 0x6f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x50, 0x61, 0x74, 0x68,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x6f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x46, 0x69,
	0x6c, 0x65, 0x50, 0x61, 0x74, 0x68, 0x22, 0x56, 0x0a, 0x06, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72,
	0x72, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72,
	0x12, 0x1c, 0x0a, 0x09, 0x73, 0x6f, 0x6d, 0x65, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x6f, 0x6d, 0x65, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x32, 0x3b,
	0x0a, 0x0c, 0x50, 0x61, 0x72, 0x73, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x2b,
	0x0a, 0x05, 0x50, 0x61, 0x72, 0x73, 0x65, 0x12, 0x12, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x50,
	0x61, 0x74, 0x68, 0x73, 0x54, 0x6f, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x1a, 0x0c, 0x2e, 0x67, 0x72,
	0x70, 0x63, 0x2e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x00, 0x42, 0x2e, 0x5a, 0x2c, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6a, 0x6f, 0x63, 0x68, 0x65, 0x6e,
	0x62, 0x6f, 0x65, 0x73, 0x6d, 0x61, 0x6e, 0x73, 0x2f, 0x67, 0x65, 0x64, 0x63, 0x6f, 0x6d, 0x2d,
	0x70, 0x61, 0x72, 0x73, 0x65, 0x72, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_grpc_parse_proto_rawDescOnce sync.Once
	file_grpc_parse_proto_rawDescData = file_grpc_parse_proto_rawDesc
)

func file_grpc_parse_proto_rawDescGZIP() []byte {
	file_grpc_parse_proto_rawDescOnce.Do(func() {
		file_grpc_parse_proto_rawDescData = protoimpl.X.CompressGZIP(file_grpc_parse_proto_rawDescData)
	})
	return file_grpc_parse_proto_rawDescData
}

var file_grpc_parse_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_grpc_parse_proto_goTypes = []interface{}{
	(*PathsToFiles)(nil), // 0: grpc.PathsToFiles
	(*Result)(nil),       // 1: grpc.Result
}
var file_grpc_parse_proto_depIdxs = []int32{
	0, // 0: grpc.ParseService.Parse:input_type -> grpc.PathsToFiles
	1, // 1: grpc.ParseService.Parse:output_type -> grpc.Result
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_grpc_parse_proto_init() }
func file_grpc_parse_proto_init() {
	if File_grpc_parse_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_grpc_parse_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PathsToFiles); i {
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
		file_grpc_parse_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Result); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_grpc_parse_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_grpc_parse_proto_goTypes,
		DependencyIndexes: file_grpc_parse_proto_depIdxs,
		MessageInfos:      file_grpc_parse_proto_msgTypes,
	}.Build()
	File_grpc_parse_proto = out.File
	file_grpc_parse_proto_rawDesc = nil
	file_grpc_parse_proto_goTypes = nil
	file_grpc_parse_proto_depIdxs = nil
}
