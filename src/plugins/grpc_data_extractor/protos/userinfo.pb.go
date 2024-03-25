// Code generated by protoc-gen-go-plugin. DO NOT EDIT.
// versions:
// 	protoc-gen-go-plugin 0.8.0
// 	protoc               v3.21.12
// source: userinfo.proto

package protos

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type PhoneType int32

const (
	PhoneType_PHONE_TYPE_UNSPECIFIED PhoneType = 0
	PhoneType_PHONE_TYPE_MOBILE      PhoneType = 1
	PhoneType_PHONE_TYPE_HOME        PhoneType = 2
	PhoneType_PHONE_TYPE_WORK        PhoneType = 3
)

// Enum value maps for PhoneType.
var (
	PhoneType_name = map[int32]string{
		0: "PHONE_TYPE_UNSPECIFIED",
		1: "PHONE_TYPE_MOBILE",
		2: "PHONE_TYPE_HOME",
		3: "PHONE_TYPE_WORK",
	}
	PhoneType_value = map[string]int32{
		"PHONE_TYPE_UNSPECIFIED": 0,
		"PHONE_TYPE_MOBILE":      1,
		"PHONE_TYPE_HOME":        2,
		"PHONE_TYPE_WORK":        3,
	}
)

func (x PhoneType) Enum() *PhoneType {
	p := new(PhoneType)
	*p = x
	return p
}

type PhoneNumber struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Number string    `protobuf:"bytes,1,opt,name=number,proto3" json:"number,omitempty"`
	Type   PhoneType `protobuf:"varint,2,opt,name=type,proto3,enum=grpc_data_extractor.PhoneType" json:"type,omitempty"`
}

func (x *PhoneNumber) ProtoReflect() protoreflect.Message {
	panic(`not implemented`)
}

func (x *PhoneNumber) GetNumber() string {
	if x != nil {
		return x.Number
	}
	return ""
}

func (x *PhoneNumber) GetType() PhoneType {
	if x != nil {
		return x.Type
	}
	return PhoneType_PHONE_TYPE_UNSPECIFIED
}

type GetUserInfoRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId int32 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *GetUserInfoRequest) ProtoReflect() protoreflect.Message {
	panic(`not implemented`)
}

func (x *GetUserInfoRequest) GetUserId() int32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

type UserInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Openid   string         `protobuf:"bytes,1,opt,name=openid,proto3" json:"openid,omitempty"`
	UserId   int32          `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Username string         `protobuf:"bytes,3,opt,name=username,proto3" json:"username,omitempty"`
	Email    string         `protobuf:"bytes,4,opt,name=email,proto3" json:"email,omitempty"`
	Phones   []*PhoneNumber `protobuf:"bytes,5,rep,name=phones,proto3" json:"phones,omitempty"`
}

func (x *UserInfo) ProtoReflect() protoreflect.Message {
	panic(`not implemented`)
}

func (x *UserInfo) GetOpenid() string {
	if x != nil {
		return x.Openid
	}
	return ""
}

func (x *UserInfo) GetUserId() int32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *UserInfo) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *UserInfo) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *UserInfo) GetPhones() []*PhoneNumber {
	if x != nil {
		return x.Phones
	}
	return nil
}
