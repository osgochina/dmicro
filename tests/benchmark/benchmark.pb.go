// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.6.1
// source: benchmark.proto

package benchmark

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

type BenchmarkMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Field1  string   `protobuf:"bytes,1,opt,name=Field1,proto3" json:"Field1,omitempty"`
	Field2  int32    `protobuf:"varint,2,opt,name=Field2,proto3" json:"Field2,omitempty"`
	Field3  int32    `protobuf:"varint,3,opt,name=Field3,proto3" json:"Field3,omitempty"`
	Field4  string   `protobuf:"bytes,4,opt,name=Field4,proto3" json:"Field4,omitempty"`
	Field5  []string `protobuf:"bytes,5,rep,name=Field5,proto3" json:"Field5,omitempty"`
	Field6  int32    `protobuf:"varint,6,opt,name=Field6,proto3" json:"Field6,omitempty"`
	Field7  string   `protobuf:"bytes,7,opt,name=Field7,proto3" json:"Field7,omitempty"`
	Field8  string   `protobuf:"bytes,8,opt,name=Field8,proto3" json:"Field8,omitempty"`
	Field9  bool     `protobuf:"varint,9,opt,name=Field9,proto3" json:"Field9,omitempty"`
	Field10 bool     `protobuf:"varint,10,opt,name=Field10,proto3" json:"Field10,omitempty"`
	Field11 bool     `protobuf:"varint,11,opt,name=Field11,proto3" json:"Field11,omitempty"`
	Field12 bool     `protobuf:"varint,12,opt,name=Field12,proto3" json:"Field12,omitempty"`
	Field13 int32    `protobuf:"varint,13,opt,name=Field13,proto3" json:"Field13,omitempty"`
	Field14 string   `protobuf:"bytes,14,opt,name=Field14,proto3" json:"Field14,omitempty"`
	Field15 int32    `protobuf:"varint,15,opt,name=Field15,proto3" json:"Field15,omitempty"`
	Field16 int64    `protobuf:"varint,16,opt,name=Field16,proto3" json:"Field16,omitempty"`
	Field17 bool     `protobuf:"varint,17,opt,name=Field17,proto3" json:"Field17,omitempty"`
	Field18 int32    `protobuf:"varint,18,opt,name=Field18,proto3" json:"Field18,omitempty"`
	Field19 int32    `protobuf:"varint,19,opt,name=Field19,proto3" json:"Field19,omitempty"`
	Field20 bool     `protobuf:"varint,20,opt,name=Field20,proto3" json:"Field20,omitempty"`
	Field21 bool     `protobuf:"varint,21,opt,name=Field21,proto3" json:"Field21,omitempty"`
	Field22 int32    `protobuf:"varint,22,opt,name=Field22,proto3" json:"Field22,omitempty"`
	Field23 int32    `protobuf:"varint,23,opt,name=Field23,proto3" json:"Field23,omitempty"`
	Field24 int32    `protobuf:"varint,24,opt,name=Field24,proto3" json:"Field24,omitempty"`
	Field25 bool     `protobuf:"varint,25,opt,name=Field25,proto3" json:"Field25,omitempty"`
	Field26 bool     `protobuf:"varint,26,opt,name=Field26,proto3" json:"Field26,omitempty"`
	Field27 bool     `protobuf:"varint,27,opt,name=Field27,proto3" json:"Field27,omitempty"`
	Field28 int32    `protobuf:"varint,28,opt,name=Field28,proto3" json:"Field28,omitempty"`
	Field29 int32    `protobuf:"varint,29,opt,name=Field29,proto3" json:"Field29,omitempty"`
	Field30 string   `protobuf:"bytes,30,opt,name=Field30,proto3" json:"Field30,omitempty"`
	Field31 string   `protobuf:"bytes,31,opt,name=Field31,proto3" json:"Field31,omitempty"`
	Field32 int32    `protobuf:"varint,32,opt,name=Field32,proto3" json:"Field32,omitempty"`
	Field33 int32    `protobuf:"varint,33,opt,name=Field33,proto3" json:"Field33,omitempty"`
	Field34 string   `protobuf:"bytes,34,opt,name=Field34,proto3" json:"Field34,omitempty"`
	Field35 int32    `protobuf:"varint,35,opt,name=Field35,proto3" json:"Field35,omitempty"`
	Field36 int32    `protobuf:"varint,36,opt,name=Field36,proto3" json:"Field36,omitempty"`
	Field37 int32    `protobuf:"varint,37,opt,name=Field37,proto3" json:"Field37,omitempty"`
	Field38 int32    `protobuf:"varint,38,opt,name=Field38,proto3" json:"Field38,omitempty"`
	Field39 int32    `protobuf:"varint,39,opt,name=Field39,proto3" json:"Field39,omitempty"`
	Field40 int32    `protobuf:"varint,40,opt,name=Field40,proto3" json:"Field40,omitempty"`
}

func (x *BenchmarkMessage) Reset() {
	*x = BenchmarkMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_benchmark_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BenchmarkMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BenchmarkMessage) ProtoMessage() {}

func (x *BenchmarkMessage) ProtoReflect() protoreflect.Message {
	mi := &file_benchmark_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BenchmarkMessage.ProtoReflect.Descriptor instead.
func (*BenchmarkMessage) Descriptor() ([]byte, []int) {
	return file_benchmark_proto_rawDescGZIP(), []int{0}
}

func (x *BenchmarkMessage) GetField1() string {
	if x != nil {
		return x.Field1
	}
	return ""
}

func (x *BenchmarkMessage) GetField2() int32 {
	if x != nil {
		return x.Field2
	}
	return 0
}

func (x *BenchmarkMessage) GetField3() int32 {
	if x != nil {
		return x.Field3
	}
	return 0
}

func (x *BenchmarkMessage) GetField4() string {
	if x != nil {
		return x.Field4
	}
	return ""
}

func (x *BenchmarkMessage) GetField5() []string {
	if x != nil {
		return x.Field5
	}
	return nil
}

func (x *BenchmarkMessage) GetField6() int32 {
	if x != nil {
		return x.Field6
	}
	return 0
}

func (x *BenchmarkMessage) GetField7() string {
	if x != nil {
		return x.Field7
	}
	return ""
}

func (x *BenchmarkMessage) GetField8() string {
	if x != nil {
		return x.Field8
	}
	return ""
}

func (x *BenchmarkMessage) GetField9() bool {
	if x != nil {
		return x.Field9
	}
	return false
}

func (x *BenchmarkMessage) GetField10() bool {
	if x != nil {
		return x.Field10
	}
	return false
}

func (x *BenchmarkMessage) GetField11() bool {
	if x != nil {
		return x.Field11
	}
	return false
}

func (x *BenchmarkMessage) GetField12() bool {
	if x != nil {
		return x.Field12
	}
	return false
}

func (x *BenchmarkMessage) GetField13() int32 {
	if x != nil {
		return x.Field13
	}
	return 0
}

func (x *BenchmarkMessage) GetField14() string {
	if x != nil {
		return x.Field14
	}
	return ""
}

func (x *BenchmarkMessage) GetField15() int32 {
	if x != nil {
		return x.Field15
	}
	return 0
}

func (x *BenchmarkMessage) GetField16() int64 {
	if x != nil {
		return x.Field16
	}
	return 0
}

func (x *BenchmarkMessage) GetField17() bool {
	if x != nil {
		return x.Field17
	}
	return false
}

func (x *BenchmarkMessage) GetField18() int32 {
	if x != nil {
		return x.Field18
	}
	return 0
}

func (x *BenchmarkMessage) GetField19() int32 {
	if x != nil {
		return x.Field19
	}
	return 0
}

func (x *BenchmarkMessage) GetField20() bool {
	if x != nil {
		return x.Field20
	}
	return false
}

func (x *BenchmarkMessage) GetField21() bool {
	if x != nil {
		return x.Field21
	}
	return false
}

func (x *BenchmarkMessage) GetField22() int32 {
	if x != nil {
		return x.Field22
	}
	return 0
}

func (x *BenchmarkMessage) GetField23() int32 {
	if x != nil {
		return x.Field23
	}
	return 0
}

func (x *BenchmarkMessage) GetField24() int32 {
	if x != nil {
		return x.Field24
	}
	return 0
}

func (x *BenchmarkMessage) GetField25() bool {
	if x != nil {
		return x.Field25
	}
	return false
}

func (x *BenchmarkMessage) GetField26() bool {
	if x != nil {
		return x.Field26
	}
	return false
}

func (x *BenchmarkMessage) GetField27() bool {
	if x != nil {
		return x.Field27
	}
	return false
}

func (x *BenchmarkMessage) GetField28() int32 {
	if x != nil {
		return x.Field28
	}
	return 0
}

func (x *BenchmarkMessage) GetField29() int32 {
	if x != nil {
		return x.Field29
	}
	return 0
}

func (x *BenchmarkMessage) GetField30() string {
	if x != nil {
		return x.Field30
	}
	return ""
}

func (x *BenchmarkMessage) GetField31() string {
	if x != nil {
		return x.Field31
	}
	return ""
}

func (x *BenchmarkMessage) GetField32() int32 {
	if x != nil {
		return x.Field32
	}
	return 0
}

func (x *BenchmarkMessage) GetField33() int32 {
	if x != nil {
		return x.Field33
	}
	return 0
}

func (x *BenchmarkMessage) GetField34() string {
	if x != nil {
		return x.Field34
	}
	return ""
}

func (x *BenchmarkMessage) GetField35() int32 {
	if x != nil {
		return x.Field35
	}
	return 0
}

func (x *BenchmarkMessage) GetField36() int32 {
	if x != nil {
		return x.Field36
	}
	return 0
}

func (x *BenchmarkMessage) GetField37() int32 {
	if x != nil {
		return x.Field37
	}
	return 0
}

func (x *BenchmarkMessage) GetField38() int32 {
	if x != nil {
		return x.Field38
	}
	return 0
}

func (x *BenchmarkMessage) GetField39() int32 {
	if x != nil {
		return x.Field39
	}
	return 0
}

func (x *BenchmarkMessage) GetField40() int32 {
	if x != nil {
		return x.Field40
	}
	return 0
}

var File_benchmark_proto protoreflect.FileDescriptor

var file_benchmark_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x62, 0x65, 0x6e, 0x63, 0x68, 0x6d, 0x61, 0x72, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x09, 0x62, 0x65, 0x6e, 0x63, 0x68, 0x6d, 0x61, 0x72, 0x6b, 0x22, 0x90, 0x08, 0x0a,
	0x10, 0x42, 0x65, 0x6e, 0x63, 0x68, 0x6d, 0x61, 0x72, 0x6b, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x12, 0x16, 0x0a, 0x06, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x31, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x31, 0x12, 0x16, 0x0a, 0x06, 0x46, 0x69, 0x65,
	0x6c, 0x64, 0x32, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x46, 0x69, 0x65, 0x6c, 0x64,
	0x32, 0x12, 0x16, 0x0a, 0x06, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x33, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x06, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x33, 0x12, 0x16, 0x0a, 0x06, 0x46, 0x69, 0x65,
	0x6c, 0x64, 0x34, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x46, 0x69, 0x65, 0x6c, 0x64,
	0x34, 0x12, 0x16, 0x0a, 0x06, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x35, 0x18, 0x05, 0x20, 0x03, 0x28,
	0x09, 0x52, 0x06, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x35, 0x12, 0x16, 0x0a, 0x06, 0x46, 0x69, 0x65,
	0x6c, 0x64, 0x36, 0x18, 0x06, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x46, 0x69, 0x65, 0x6c, 0x64,
	0x36, 0x12, 0x16, 0x0a, 0x06, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x37, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x37, 0x12, 0x16, 0x0a, 0x06, 0x46, 0x69, 0x65,
	0x6c, 0x64, 0x38, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x46, 0x69, 0x65, 0x6c, 0x64,
	0x38, 0x12, 0x16, 0x0a, 0x06, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x39, 0x18, 0x09, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x06, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x39, 0x12, 0x18, 0x0a, 0x07, 0x46, 0x69, 0x65,
	0x6c, 0x64, 0x31, 0x30, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x46, 0x69, 0x65, 0x6c,
	0x64, 0x31, 0x30, 0x12, 0x18, 0x0a, 0x07, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x31, 0x31, 0x18, 0x0b,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x31, 0x31, 0x12, 0x18, 0x0a,
	0x07, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x31, 0x32, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07,
	0x46, 0x69, 0x65, 0x6c, 0x64, 0x31, 0x32, 0x12, 0x18, 0x0a, 0x07, 0x46, 0x69, 0x65, 0x6c, 0x64,
	0x31, 0x33, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x31,
	0x33, 0x12, 0x18, 0x0a, 0x07, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x31, 0x34, 0x18, 0x0e, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x31, 0x34, 0x12, 0x18, 0x0a, 0x07, 0x46,
	0x69, 0x65, 0x6c, 0x64, 0x31, 0x35, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x46, 0x69,
	0x65, 0x6c, 0x64, 0x31, 0x35, 0x12, 0x18, 0x0a, 0x07, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x31, 0x36,
	0x18, 0x10, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x31, 0x36, 0x12,
	0x18, 0x0a, 0x07, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x31, 0x37, 0x18, 0x11, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x07, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x31, 0x37, 0x12, 0x18, 0x0a, 0x07, 0x46, 0x69, 0x65,
	0x6c, 0x64, 0x31, 0x38, 0x18, 0x12, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x46, 0x69, 0x65, 0x6c,
	0x64, 0x31, 0x38, 0x12, 0x18, 0x0a, 0x07, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x31, 0x39, 0x18, 0x13,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x31, 0x39, 0x12, 0x18, 0x0a,
	0x07, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x32, 0x30, 0x18, 0x14, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07,
	0x46, 0x69, 0x65, 0x6c, 0x64, 0x32, 0x30, 0x12, 0x18, 0x0a, 0x07, 0x46, 0x69, 0x65, 0x6c, 0x64,
	0x32, 0x31, 0x18, 0x15, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x32,
	0x31, 0x12, 0x18, 0x0a, 0x07, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x32, 0x32, 0x18, 0x16, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x07, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x32, 0x32, 0x12, 0x18, 0x0a, 0x07, 0x46,
	0x69, 0x65, 0x6c, 0x64, 0x32, 0x33, 0x18, 0x17, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x46, 0x69,
	0x65, 0x6c, 0x64, 0x32, 0x33, 0x12, 0x18, 0x0a, 0x07, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x32, 0x34,
	0x18, 0x18, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x32, 0x34, 0x12,
	0x18, 0x0a, 0x07, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x32, 0x35, 0x18, 0x19, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x07, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x32, 0x35, 0x12, 0x18, 0x0a, 0x07, 0x46, 0x69, 0x65,
	0x6c, 0x64, 0x32, 0x36, 0x18, 0x1a, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x46, 0x69, 0x65, 0x6c,
	0x64, 0x32, 0x36, 0x12, 0x18, 0x0a, 0x07, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x32, 0x37, 0x18, 0x1b,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x32, 0x37, 0x12, 0x18, 0x0a,
	0x07, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x32, 0x38, 0x18, 0x1c, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07,
	0x46, 0x69, 0x65, 0x6c, 0x64, 0x32, 0x38, 0x12, 0x18, 0x0a, 0x07, 0x46, 0x69, 0x65, 0x6c, 0x64,
	0x32, 0x39, 0x18, 0x1d, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x32,
	0x39, 0x12, 0x18, 0x0a, 0x07, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x33, 0x30, 0x18, 0x1e, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x33, 0x30, 0x12, 0x18, 0x0a, 0x07, 0x46,
	0x69, 0x65, 0x6c, 0x64, 0x33, 0x31, 0x18, 0x1f, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x46, 0x69,
	0x65, 0x6c, 0x64, 0x33, 0x31, 0x12, 0x18, 0x0a, 0x07, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x33, 0x32,
	0x18, 0x20, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x33, 0x32, 0x12,
	0x18, 0x0a, 0x07, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x33, 0x33, 0x18, 0x21, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x07, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x33, 0x33, 0x12, 0x18, 0x0a, 0x07, 0x46, 0x69, 0x65,
	0x6c, 0x64, 0x33, 0x34, 0x18, 0x22, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x46, 0x69, 0x65, 0x6c,
	0x64, 0x33, 0x34, 0x12, 0x18, 0x0a, 0x07, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x33, 0x35, 0x18, 0x23,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x33, 0x35, 0x12, 0x18, 0x0a,
	0x07, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x33, 0x36, 0x18, 0x24, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07,
	0x46, 0x69, 0x65, 0x6c, 0x64, 0x33, 0x36, 0x12, 0x18, 0x0a, 0x07, 0x46, 0x69, 0x65, 0x6c, 0x64,
	0x33, 0x37, 0x18, 0x25, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x33,
	0x37, 0x12, 0x18, 0x0a, 0x07, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x33, 0x38, 0x18, 0x26, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x07, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x33, 0x38, 0x12, 0x18, 0x0a, 0x07, 0x46,
	0x69, 0x65, 0x6c, 0x64, 0x33, 0x39, 0x18, 0x27, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x46, 0x69,
	0x65, 0x6c, 0x64, 0x33, 0x39, 0x12, 0x18, 0x0a, 0x07, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x34, 0x30,
	0x18, 0x28, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x34, 0x30, 0x42,
	0x10, 0x48, 0x01, 0x5a, 0x0c, 0x2e, 0x2f, 0x3b, 0x62, 0x65, 0x6e, 0x63, 0x68, 0x6d, 0x61, 0x72,
	0x6b, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_benchmark_proto_rawDescOnce sync.Once
	file_benchmark_proto_rawDescData = file_benchmark_proto_rawDesc
)

func file_benchmark_proto_rawDescGZIP() []byte {
	file_benchmark_proto_rawDescOnce.Do(func() {
		file_benchmark_proto_rawDescData = protoimpl.X.CompressGZIP(file_benchmark_proto_rawDescData)
	})
	return file_benchmark_proto_rawDescData
}

var file_benchmark_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_benchmark_proto_goTypes = []interface{}{
	(*BenchmarkMessage)(nil), // 0: benchmark.BenchmarkMessage
}
var file_benchmark_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_benchmark_proto_init() }
func file_benchmark_proto_init() {
	if File_benchmark_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_benchmark_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BenchmarkMessage); i {
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
			RawDescriptor: file_benchmark_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_benchmark_proto_goTypes,
		DependencyIndexes: file_benchmark_proto_depIdxs,
		MessageInfos:      file_benchmark_proto_msgTypes,
	}.Build()
	File_benchmark_proto = out.File
	file_benchmark_proto_rawDesc = nil
	file_benchmark_proto_goTypes = nil
	file_benchmark_proto_depIdxs = nil
}