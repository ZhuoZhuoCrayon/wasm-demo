// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: userinfo.proto

package userinfoservice

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

func (x PhoneType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PhoneType) Descriptor() protoreflect.EnumDescriptor {
	return file_userinfo_proto_enumTypes[0].Descriptor()
}

func (PhoneType) Type() protoreflect.EnumType {
	return &file_userinfo_proto_enumTypes[0]
}

func (x PhoneType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PhoneType.Descriptor instead.
func (PhoneType) EnumDescriptor() ([]byte, []int) {
	return file_userinfo_proto_rawDescGZIP(), []int{0}
}

type PhoneNumber struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Number string    `protobuf:"bytes,1,opt,name=number,proto3" json:"number,omitempty"`
	Type   PhoneType `protobuf:"varint,2,opt,name=type,proto3,enum=userinfoservice.PhoneType" json:"type,omitempty"`
}

func (x *PhoneNumber) Reset() {
	*x = PhoneNumber{}
	if protoimpl.UnsafeEnabled {
		mi := &file_userinfo_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PhoneNumber) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PhoneNumber) ProtoMessage() {}

func (x *PhoneNumber) ProtoReflect() protoreflect.Message {
	mi := &file_userinfo_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PhoneNumber.ProtoReflect.Descriptor instead.
func (*PhoneNumber) Descriptor() ([]byte, []int) {
	return file_userinfo_proto_rawDescGZIP(), []int{0}
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

func (x *GetUserInfoRequest) Reset() {
	*x = GetUserInfoRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_userinfo_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUserInfoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserInfoRequest) ProtoMessage() {}

func (x *GetUserInfoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_userinfo_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserInfoRequest.ProtoReflect.Descriptor instead.
func (*GetUserInfoRequest) Descriptor() ([]byte, []int) {
	return file_userinfo_proto_rawDescGZIP(), []int{1}
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

func (x *UserInfo) Reset() {
	*x = UserInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_userinfo_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserInfo) ProtoMessage() {}

func (x *UserInfo) ProtoReflect() protoreflect.Message {
	mi := &file_userinfo_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserInfo.ProtoReflect.Descriptor instead.
func (*UserInfo) Descriptor() ([]byte, []int) {
	return file_userinfo_proto_rawDescGZIP(), []int{2}
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

var File_userinfo_proto protoreflect.FileDescriptor

var file_userinfo_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x75, 0x73, 0x65, 0x72, 0x69, 0x6e, 0x66, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x0f, 0x75, 0x73, 0x65, 0x72, 0x69, 0x6e, 0x66, 0x6f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x22, 0x55, 0x0a, 0x0b, 0x50, 0x68, 0x6f, 0x6e, 0x65, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72,
	0x12, 0x16, 0x0a, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x2e, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1a, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x69, 0x6e, 0x66,
	0x6f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x50, 0x68, 0x6f, 0x6e, 0x65, 0x54, 0x79,
	0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x2d, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x55,
	0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17,
	0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0xa3, 0x01, 0x0a, 0x08, 0x55, 0x73, 0x65, 0x72,
	0x49, 0x6e, 0x66, 0x6f, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x70, 0x65, 0x6e, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6f, 0x70, 0x65, 0x6e, 0x69, 0x64, 0x12, 0x17, 0x0a, 0x07,
	0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x75,
	0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x34, 0x0a, 0x06, 0x70, 0x68, 0x6f, 0x6e, 0x65,
	0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x69, 0x6e,
	0x66, 0x6f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x50, 0x68, 0x6f, 0x6e, 0x65, 0x4e,
	0x75, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x06, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x73, 0x2a, 0x68, 0x0a,
	0x09, 0x50, 0x68, 0x6f, 0x6e, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1a, 0x0a, 0x16, 0x50, 0x48,
	0x4f, 0x4e, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49,
	0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x15, 0x0a, 0x11, 0x50, 0x48, 0x4f, 0x4e, 0x45, 0x5f,
	0x54, 0x59, 0x50, 0x45, 0x5f, 0x4d, 0x4f, 0x42, 0x49, 0x4c, 0x45, 0x10, 0x01, 0x12, 0x13, 0x0a,
	0x0f, 0x50, 0x48, 0x4f, 0x4e, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x48, 0x4f, 0x4d, 0x45,
	0x10, 0x02, 0x12, 0x13, 0x0a, 0x0f, 0x50, 0x48, 0x4f, 0x4e, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45,
	0x5f, 0x57, 0x4f, 0x52, 0x4b, 0x10, 0x03, 0x32, 0x60, 0x0a, 0x0f, 0x55, 0x73, 0x65, 0x72, 0x49,
	0x6e, 0x66, 0x6f, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x4d, 0x0a, 0x0b, 0x47, 0x65,
	0x74, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x23, 0x2e, 0x75, 0x73, 0x65, 0x72,
	0x69, 0x6e, 0x66, 0x6f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x55,
	0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19,
	0x2e, 0x75, 0x73, 0x65, 0x72, 0x69, 0x6e, 0x66, 0x6f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x2e, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x42, 0x39, 0x5a, 0x37, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x5a, 0x68, 0x75, 0x6f, 0x5a, 0x68, 0x75, 0x6f,
	0x43, 0x72, 0x61, 0x79, 0x6f, 0x6e, 0x2f, 0x77, 0x61, 0x73, 0x6d, 0x2d, 0x64, 0x65, 0x6d, 0x6f,
	0x2f, 0x73, 0x72, 0x63, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x69, 0x6e, 0x66, 0x6f, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_userinfo_proto_rawDescOnce sync.Once
	file_userinfo_proto_rawDescData = file_userinfo_proto_rawDesc
)

func file_userinfo_proto_rawDescGZIP() []byte {
	file_userinfo_proto_rawDescOnce.Do(func() {
		file_userinfo_proto_rawDescData = protoimpl.X.CompressGZIP(file_userinfo_proto_rawDescData)
	})
	return file_userinfo_proto_rawDescData
}

var file_userinfo_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_userinfo_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_userinfo_proto_goTypes = []interface{}{
	(PhoneType)(0),             // 0: userinfoservice.PhoneType
	(*PhoneNumber)(nil),        // 1: userinfoservice.PhoneNumber
	(*GetUserInfoRequest)(nil), // 2: userinfoservice.GetUserInfoRequest
	(*UserInfo)(nil),           // 3: userinfoservice.UserInfo
}
var file_userinfo_proto_depIdxs = []int32{
	0, // 0: userinfoservice.PhoneNumber.type:type_name -> userinfoservice.PhoneType
	1, // 1: userinfoservice.UserInfo.phones:type_name -> userinfoservice.PhoneNumber
	2, // 2: userinfoservice.UserInfoService.GetUserInfo:input_type -> userinfoservice.GetUserInfoRequest
	3, // 3: userinfoservice.UserInfoService.GetUserInfo:output_type -> userinfoservice.UserInfo
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_userinfo_proto_init() }
func file_userinfo_proto_init() {
	if File_userinfo_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_userinfo_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PhoneNumber); i {
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
		file_userinfo_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUserInfoRequest); i {
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
		file_userinfo_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserInfo); i {
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
			RawDescriptor: file_userinfo_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_userinfo_proto_goTypes,
		DependencyIndexes: file_userinfo_proto_depIdxs,
		EnumInfos:         file_userinfo_proto_enumTypes,
		MessageInfos:      file_userinfo_proto_msgTypes,
	}.Build()
	File_userinfo_proto = out.File
	file_userinfo_proto_rawDesc = nil
	file_userinfo_proto_goTypes = nil
	file_userinfo_proto_depIdxs = nil
}