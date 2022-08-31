// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.3.0
// source: fatima.message.v1.proto

package fatima_message_v1

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

type ResponseError_GrpcResponse int32

const (
	ResponseError_UNIVERSAL           ResponseError_GrpcResponse = 0
	ResponseError_BAD_PARAMETER       ResponseError_GrpcResponse = 400
	ResponseError_UNAUTORIZED         ResponseError_GrpcResponse = 401
	ResponseError_FORBIDDEN           ResponseError_GrpcResponse = 403
	ResponseError_NOT_FOUND           ResponseError_GrpcResponse = 404
	ResponseError_NOT_ACCEPTABLE      ResponseError_GrpcResponse = 406
	ResponseError_SERVER_ERROR        ResponseError_GrpcResponse = 500
	ResponseError_SERVICE_UNAVAILABLE ResponseError_GrpcResponse = 503
)

// Enum value maps for ResponseError_GrpcResponse.
var (
	ResponseError_GrpcResponse_name = map[int32]string{
		0:   "UNIVERSAL",
		400: "BAD_PARAMETER",
		401: "UNAUTORIZED",
		403: "FORBIDDEN",
		404: "NOT_FOUND",
		406: "NOT_ACCEPTABLE",
		500: "SERVER_ERROR",
		503: "SERVICE_UNAVAILABLE",
	}
	ResponseError_GrpcResponse_value = map[string]int32{
		"UNIVERSAL":           0,
		"BAD_PARAMETER":       400,
		"UNAUTORIZED":         401,
		"FORBIDDEN":           403,
		"NOT_FOUND":           404,
		"NOT_ACCEPTABLE":      406,
		"SERVER_ERROR":        500,
		"SERVICE_UNAVAILABLE": 503,
	}
)

func (x ResponseError_GrpcResponse) Enum() *ResponseError_GrpcResponse {
	p := new(ResponseError_GrpcResponse)
	*p = x
	return p
}

func (x ResponseError_GrpcResponse) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ResponseError_GrpcResponse) Descriptor() protoreflect.EnumDescriptor {
	return file_fatima_message_v1_proto_enumTypes[0].Descriptor()
}

func (ResponseError_GrpcResponse) Type() protoreflect.EnumType {
	return &file_fatima_message_v1_proto_enumTypes[0]
}

func (x ResponseError_GrpcResponse) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ResponseError_GrpcResponse.Descriptor instead.
func (ResponseError_GrpcResponse) EnumDescriptor() ([]byte, []int) {
	return file_fatima_message_v1_proto_rawDescGZIP(), []int{3, 0}
}

type ResponseError_ErrorCode int32

const (
	ResponseError_SUCCESS        ResponseError_ErrorCode = 0
	ResponseError_NO_RECORD      ResponseError_ErrorCode = 100
	ResponseError_ERROR_RESPONSE ResponseError_ErrorCode = 101
	ResponseError_CONNECT_FAIL   ResponseError_ErrorCode = 102
	ResponseError_ERROR_ETC      ResponseError_ErrorCode = 103
)

// Enum value maps for ResponseError_ErrorCode.
var (
	ResponseError_ErrorCode_name = map[int32]string{
		0:   "SUCCESS",
		100: "NO_RECORD",
		101: "ERROR_RESPONSE",
		102: "CONNECT_FAIL",
		103: "ERROR_ETC",
	}
	ResponseError_ErrorCode_value = map[string]int32{
		"SUCCESS":        0,
		"NO_RECORD":      100,
		"ERROR_RESPONSE": 101,
		"CONNECT_FAIL":   102,
		"ERROR_ETC":      103,
	}
)

func (x ResponseError_ErrorCode) Enum() *ResponseError_ErrorCode {
	p := new(ResponseError_ErrorCode)
	*p = x
	return p
}

func (x ResponseError_ErrorCode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ResponseError_ErrorCode) Descriptor() protoreflect.EnumDescriptor {
	return file_fatima_message_v1_proto_enumTypes[1].Descriptor()
}

func (ResponseError_ErrorCode) Type() protoreflect.EnumType {
	return &file_fatima_message_v1_proto_enumTypes[1]
}

func (x ResponseError_ErrorCode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ResponseError_ErrorCode.Descriptor instead.
func (ResponseError_ErrorCode) EnumDescriptor() ([]byte, []int) {
	return file_fatima_message_v1_proto_rawDescGZIP(), []int{3, 1}
}

type SendFatimaMessageRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	JsonString string `protobuf:"bytes,1,opt,name=jsonString,proto3" json:"jsonString,omitempty"`
}

func (x *SendFatimaMessageRequest) Reset() {
	*x = SendFatimaMessageRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fatima_message_v1_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendFatimaMessageRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendFatimaMessageRequest) ProtoMessage() {}

func (x *SendFatimaMessageRequest) ProtoReflect() protoreflect.Message {
	mi := &file_fatima_message_v1_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendFatimaMessageRequest.ProtoReflect.Descriptor instead.
func (*SendFatimaMessageRequest) Descriptor() ([]byte, []int) {
	return file_fatima_message_v1_proto_rawDescGZIP(), []int{0}
}

func (x *SendFatimaMessageRequest) GetJsonString() string {
	if x != nil {
		return x.JsonString
	}
	return ""
}

type SendFatimaMessageResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Response:
	//	*SendFatimaMessageResponse_Success
	//	*SendFatimaMessageResponse_Error
	Response isSendFatimaMessageResponse_Response `protobuf_oneof:"response"`
}

func (x *SendFatimaMessageResponse) Reset() {
	*x = SendFatimaMessageResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fatima_message_v1_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendFatimaMessageResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendFatimaMessageResponse) ProtoMessage() {}

func (x *SendFatimaMessageResponse) ProtoReflect() protoreflect.Message {
	mi := &file_fatima_message_v1_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendFatimaMessageResponse.ProtoReflect.Descriptor instead.
func (*SendFatimaMessageResponse) Descriptor() ([]byte, []int) {
	return file_fatima_message_v1_proto_rawDescGZIP(), []int{1}
}

func (m *SendFatimaMessageResponse) GetResponse() isSendFatimaMessageResponse_Response {
	if m != nil {
		return m.Response
	}
	return nil
}

func (x *SendFatimaMessageResponse) GetSuccess() *ResponseSuccess {
	if x, ok := x.GetResponse().(*SendFatimaMessageResponse_Success); ok {
		return x.Success
	}
	return nil
}

func (x *SendFatimaMessageResponse) GetError() *ResponseError {
	if x, ok := x.GetResponse().(*SendFatimaMessageResponse_Error); ok {
		return x.Error
	}
	return nil
}

type isSendFatimaMessageResponse_Response interface {
	isSendFatimaMessageResponse_Response()
}

type SendFatimaMessageResponse_Success struct {
	Success *ResponseSuccess `protobuf:"bytes,1,opt,name=success,proto3,oneof"`
}

type SendFatimaMessageResponse_Error struct {
	Error *ResponseError `protobuf:"bytes,2,opt,name=error,proto3,oneof"`
}

func (*SendFatimaMessageResponse_Success) isSendFatimaMessageResponse_Response() {}

func (*SendFatimaMessageResponse_Error) isSendFatimaMessageResponse_Response() {}

type ResponseSuccess struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ResponseSuccess) Reset() {
	*x = ResponseSuccess{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fatima_message_v1_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResponseSuccess) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResponseSuccess) ProtoMessage() {}

func (x *ResponseSuccess) ProtoReflect() protoreflect.Message {
	mi := &file_fatima_message_v1_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResponseSuccess.ProtoReflect.Descriptor instead.
func (*ResponseSuccess) Descriptor() ([]byte, []int) {
	return file_fatima_message_v1_proto_rawDescGZIP(), []int{2}
}

type ResponseError struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GrpcResponse ResponseError_GrpcResponse `protobuf:"varint,1,opt,name=grpcResponse,proto3,enum=fatima.message.v1.ResponseError_GrpcResponse" json:"grpcResponse,omitempty"`
	Code         ResponseError_ErrorCode    `protobuf:"varint,2,opt,name=code,proto3,enum=fatima.message.v1.ResponseError_ErrorCode" json:"code,omitempty"`
	Value        string                     `protobuf:"bytes,3,opt,name=value,proto3" json:"value,omitempty"`
	Desc         string                     `protobuf:"bytes,4,opt,name=desc,proto3" json:"desc,omitempty"`
}

func (x *ResponseError) Reset() {
	*x = ResponseError{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fatima_message_v1_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResponseError) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResponseError) ProtoMessage() {}

func (x *ResponseError) ProtoReflect() protoreflect.Message {
	mi := &file_fatima_message_v1_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResponseError.ProtoReflect.Descriptor instead.
func (*ResponseError) Descriptor() ([]byte, []int) {
	return file_fatima_message_v1_proto_rawDescGZIP(), []int{3}
}

func (x *ResponseError) GetGrpcResponse() ResponseError_GrpcResponse {
	if x != nil {
		return x.GrpcResponse
	}
	return ResponseError_UNIVERSAL
}

func (x *ResponseError) GetCode() ResponseError_ErrorCode {
	if x != nil {
		return x.Code
	}
	return ResponseError_SUCCESS
}

func (x *ResponseError) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

func (x *ResponseError) GetDesc() string {
	if x != nil {
		return x.Desc
	}
	return ""
}

var File_fatima_message_v1_proto protoreflect.FileDescriptor

var file_fatima_message_v1_proto_rawDesc = []byte{
	0x0a, 0x17, 0x66, 0x61, 0x74, 0x69, 0x6d, 0x61, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x2e, 0x76, 0x31, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x11, 0x66, 0x61, 0x74, 0x69, 0x6d,
	0x61, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x22, 0x3a, 0x0a, 0x18,
	0x53, 0x65, 0x6e, 0x64, 0x46, 0x61, 0x74, 0x69, 0x6d, 0x61, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x6a, 0x73, 0x6f, 0x6e,
	0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x6a, 0x73,
	0x6f, 0x6e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x22, 0xa1, 0x01, 0x0a, 0x19, 0x53, 0x65, 0x6e,
	0x64, 0x46, 0x61, 0x74, 0x69, 0x6d, 0x61, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3e, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x66, 0x61, 0x74, 0x69, 0x6d, 0x61,
	0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x48, 0x00, 0x52, 0x07, 0x73,
	0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x38, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x66, 0x61, 0x74, 0x69, 0x6d, 0x61, 0x2e, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x48, 0x00, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72,
	0x42, 0x0a, 0x0a, 0x08, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x11, 0x0a, 0x0f,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x22,
	0xd2, 0x03, 0x0a, 0x0d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x45, 0x72, 0x72, 0x6f,
	0x72, 0x12, 0x51, 0x0a, 0x0c, 0x67, 0x72, 0x70, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x2d, 0x2e, 0x66, 0x61, 0x74, 0x69, 0x6d, 0x61,
	0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x2e, 0x47, 0x72, 0x70, 0x63, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x52, 0x0c, 0x67, 0x72, 0x70, 0x63, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3e, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x2a, 0x2e, 0x66, 0x61, 0x74, 0x69, 0x6d, 0x61, 0x2e, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x45,
	0x72, 0x72, 0x6f, 0x72, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x04,
	0x63, 0x6f, 0x64, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x65,
	0x73, 0x63, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x64, 0x65, 0x73, 0x63, 0x22, 0xa5,
	0x01, 0x0a, 0x0c, 0x47, 0x72, 0x70, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x0d, 0x0a, 0x09, 0x55, 0x4e, 0x49, 0x56, 0x45, 0x52, 0x53, 0x41, 0x4c, 0x10, 0x00, 0x12, 0x12,
	0x0a, 0x0d, 0x42, 0x41, 0x44, 0x5f, 0x50, 0x41, 0x52, 0x41, 0x4d, 0x45, 0x54, 0x45, 0x52, 0x10,
	0x90, 0x03, 0x12, 0x10, 0x0a, 0x0b, 0x55, 0x4e, 0x41, 0x55, 0x54, 0x4f, 0x52, 0x49, 0x5a, 0x45,
	0x44, 0x10, 0x91, 0x03, 0x12, 0x0e, 0x0a, 0x09, 0x46, 0x4f, 0x52, 0x42, 0x49, 0x44, 0x44, 0x45,
	0x4e, 0x10, 0x93, 0x03, 0x12, 0x0e, 0x0a, 0x09, 0x4e, 0x4f, 0x54, 0x5f, 0x46, 0x4f, 0x55, 0x4e,
	0x44, 0x10, 0x94, 0x03, 0x12, 0x13, 0x0a, 0x0e, 0x4e, 0x4f, 0x54, 0x5f, 0x41, 0x43, 0x43, 0x45,
	0x50, 0x54, 0x41, 0x42, 0x4c, 0x45, 0x10, 0x96, 0x03, 0x12, 0x11, 0x0a, 0x0c, 0x53, 0x45, 0x52,
	0x56, 0x45, 0x52, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0xf4, 0x03, 0x12, 0x18, 0x0a, 0x13,
	0x53, 0x45, 0x52, 0x56, 0x49, 0x43, 0x45, 0x5f, 0x55, 0x4e, 0x41, 0x56, 0x41, 0x49, 0x4c, 0x41,
	0x42, 0x4c, 0x45, 0x10, 0xf7, 0x03, 0x22, 0x5c, 0x0a, 0x09, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x43,
	0x6f, 0x64, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x53, 0x55, 0x43, 0x43, 0x45, 0x53, 0x53, 0x10, 0x00,
	0x12, 0x0d, 0x0a, 0x09, 0x4e, 0x4f, 0x5f, 0x52, 0x45, 0x43, 0x4f, 0x52, 0x44, 0x10, 0x64, 0x12,
	0x12, 0x0a, 0x0e, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x5f, 0x52, 0x45, 0x53, 0x50, 0x4f, 0x4e, 0x53,
	0x45, 0x10, 0x65, 0x12, 0x10, 0x0a, 0x0c, 0x43, 0x4f, 0x4e, 0x4e, 0x45, 0x43, 0x54, 0x5f, 0x46,
	0x41, 0x49, 0x4c, 0x10, 0x66, 0x12, 0x0d, 0x0a, 0x09, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x5f, 0x45,
	0x54, 0x43, 0x10, 0x67, 0x32, 0x88, 0x01, 0x0a, 0x14, 0x46, 0x61, 0x74, 0x69, 0x6d, 0x61, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x70, 0x0a,
	0x11, 0x53, 0x65, 0x6e, 0x64, 0x46, 0x61, 0x74, 0x69, 0x6d, 0x61, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x12, 0x2b, 0x2e, 0x66, 0x61, 0x74, 0x69, 0x6d, 0x61, 0x2e, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x46, 0x61, 0x74, 0x69, 0x6d,
	0x61, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x2c, 0x2e, 0x66, 0x61, 0x74, 0x69, 0x6d, 0x61, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x46, 0x61, 0x74, 0x69, 0x6d, 0x61, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42,
	0x15, 0x5a, 0x13, 0x2e, 0x3b, 0x66, 0x61, 0x74, 0x69, 0x6d, 0x61, 0x5f, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x5f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_fatima_message_v1_proto_rawDescOnce sync.Once
	file_fatima_message_v1_proto_rawDescData = file_fatima_message_v1_proto_rawDesc
)

func file_fatima_message_v1_proto_rawDescGZIP() []byte {
	file_fatima_message_v1_proto_rawDescOnce.Do(func() {
		file_fatima_message_v1_proto_rawDescData = protoimpl.X.CompressGZIP(file_fatima_message_v1_proto_rawDescData)
	})
	return file_fatima_message_v1_proto_rawDescData
}

var file_fatima_message_v1_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_fatima_message_v1_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_fatima_message_v1_proto_goTypes = []interface{}{
	(ResponseError_GrpcResponse)(0),   // 0: fatima.message.v1.ResponseError.GrpcResponse
	(ResponseError_ErrorCode)(0),      // 1: fatima.message.v1.ResponseError.ErrorCode
	(*SendFatimaMessageRequest)(nil),  // 2: fatima.message.v1.SendFatimaMessageRequest
	(*SendFatimaMessageResponse)(nil), // 3: fatima.message.v1.SendFatimaMessageResponse
	(*ResponseSuccess)(nil),           // 4: fatima.message.v1.ResponseSuccess
	(*ResponseError)(nil),             // 5: fatima.message.v1.ResponseError
}
var file_fatima_message_v1_proto_depIdxs = []int32{
	4, // 0: fatima.message.v1.SendFatimaMessageResponse.success:type_name -> fatima.message.v1.ResponseSuccess
	5, // 1: fatima.message.v1.SendFatimaMessageResponse.error:type_name -> fatima.message.v1.ResponseError
	0, // 2: fatima.message.v1.ResponseError.grpcResponse:type_name -> fatima.message.v1.ResponseError.GrpcResponse
	1, // 3: fatima.message.v1.ResponseError.code:type_name -> fatima.message.v1.ResponseError.ErrorCode
	2, // 4: fatima.message.v1.FatimaMessageService.SendFatimaMessage:input_type -> fatima.message.v1.SendFatimaMessageRequest
	3, // 5: fatima.message.v1.FatimaMessageService.SendFatimaMessage:output_type -> fatima.message.v1.SendFatimaMessageResponse
	5, // [5:6] is the sub-list for method output_type
	4, // [4:5] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_fatima_message_v1_proto_init() }
func file_fatima_message_v1_proto_init() {
	if File_fatima_message_v1_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_fatima_message_v1_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendFatimaMessageRequest); i {
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
		file_fatima_message_v1_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendFatimaMessageResponse); i {
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
		file_fatima_message_v1_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResponseSuccess); i {
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
		file_fatima_message_v1_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResponseError); i {
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
	file_fatima_message_v1_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*SendFatimaMessageResponse_Success)(nil),
		(*SendFatimaMessageResponse_Error)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_fatima_message_v1_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_fatima_message_v1_proto_goTypes,
		DependencyIndexes: file_fatima_message_v1_proto_depIdxs,
		EnumInfos:         file_fatima_message_v1_proto_enumTypes,
		MessageInfos:      file_fatima_message_v1_proto_msgTypes,
	}.Build()
	File_fatima_message_v1_proto = out.File
	file_fatima_message_v1_proto_rawDesc = nil
	file_fatima_message_v1_proto_goTypes = nil
	file_fatima_message_v1_proto_depIdxs = nil
}