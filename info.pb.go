// Code generated by protoc-gen-go. DO NOT EDIT.
// source: info.proto

package main

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// The request message containing the user's name.
type InfoRequest struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Age                  int64    `protobuf:"varint,2,opt,name=age,proto3" json:"age,omitempty"`
	Height               int64    `protobuf:"varint,3,opt,name=height,proto3" json:"height,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InfoRequest) Reset()         { *m = InfoRequest{} }
func (m *InfoRequest) String() string { return proto.CompactTextString(m) }
func (*InfoRequest) ProtoMessage()    {}
func (*InfoRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_f140d5b28dddb141, []int{0}
}

func (m *InfoRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InfoRequest.Unmarshal(m, b)
}
func (m *InfoRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InfoRequest.Marshal(b, m, deterministic)
}
func (m *InfoRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InfoRequest.Merge(m, src)
}
func (m *InfoRequest) XXX_Size() int {
	return xxx_messageInfo_InfoRequest.Size(m)
}
func (m *InfoRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_InfoRequest.DiscardUnknown(m)
}

var xxx_messageInfo_InfoRequest proto.InternalMessageInfo

func (m *InfoRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *InfoRequest) GetAge() int64 {
	if m != nil {
		return m.Age
	}
	return 0
}

func (m *InfoRequest) GetHeight() int64 {
	if m != nil {
		return m.Height
	}
	return 0
}

// The response message containing the greetings
type InfoReply struct {
	Success              bool     `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Reason               string   `protobuf:"bytes,2,opt,name=reason,proto3" json:"reason,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InfoReply) Reset()         { *m = InfoReply{} }
func (m *InfoReply) String() string { return proto.CompactTextString(m) }
func (*InfoReply) ProtoMessage()    {}
func (*InfoReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_f140d5b28dddb141, []int{1}
}

func (m *InfoReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InfoReply.Unmarshal(m, b)
}
func (m *InfoReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InfoReply.Marshal(b, m, deterministic)
}
func (m *InfoReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InfoReply.Merge(m, src)
}
func (m *InfoReply) XXX_Size() int {
	return xxx_messageInfo_InfoReply.Size(m)
}
func (m *InfoReply) XXX_DiscardUnknown() {
	xxx_messageInfo_InfoReply.DiscardUnknown(m)
}

var xxx_messageInfo_InfoReply proto.InternalMessageInfo

func (m *InfoReply) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *InfoReply) GetReason() string {
	if m != nil {
		return m.Reason
	}
	return ""
}

func init() {
	proto.RegisterType((*InfoRequest)(nil), "main.InfoRequest")
	proto.RegisterType((*InfoReply)(nil), "main.InfoReply")
}

func init() { proto.RegisterFile("info.proto", fileDescriptor_f140d5b28dddb141) }

var fileDescriptor_f140d5b28dddb141 = []byte{
	// 206 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x90, 0xcd, 0x4e, 0x84, 0x30,
	0x14, 0x85, 0xad, 0x10, 0x90, 0x6b, 0x8c, 0x7a, 0x17, 0xa6, 0x71, 0x45, 0x58, 0xb1, 0x42, 0xa3,
	0x2b, 0x17, 0x3e, 0x80, 0x71, 0x57, 0x9e, 0xa0, 0x92, 0xcb, 0x4f, 0x02, 0x2d, 0xb6, 0x65, 0x32,
	0xbc, 0xfd, 0xa4, 0x1d, 0x26, 0x61, 0x39, 0xbb, 0x7b, 0x4e, 0xf2, 0x7d, 0x39, 0xb9, 0x00, 0x83,
	0x6a, 0x75, 0x35, 0x1b, 0xed, 0x34, 0xc6, 0x93, 0x1c, 0x54, 0xf1, 0x0b, 0xf7, 0x3f, 0xaa, 0xd5,
	0x82, 0xfe, 0x17, 0xb2, 0x0e, 0x11, 0x62, 0x25, 0x27, 0xe2, 0x2c, 0x67, 0x65, 0x26, 0xc2, 0x8d,
	0x4f, 0x10, 0xc9, 0x8e, 0xf8, 0x6d, 0xce, 0xca, 0x48, 0xf8, 0x13, 0x5f, 0x20, 0xe9, 0x69, 0xe8,
	0x7a, 0xc7, 0xa3, 0x50, 0x6e, 0xa9, 0xf8, 0x86, 0xec, 0x2c, 0x9b, 0xc7, 0x15, 0x39, 0xa4, 0x76,
	0x69, 0x1a, 0xb2, 0x36, 0xd8, 0xee, 0xc4, 0x25, 0x7a, 0xdc, 0x90, 0xb4, 0x5a, 0x05, 0x67, 0x26,
	0xb6, 0xf4, 0x71, 0x04, 0xf0, 0x78, 0x4d, 0xe6, 0x40, 0x06, 0xdf, 0x20, 0xad, 0xc9, 0xf9, 0x02,
	0x9f, 0x2b, 0xbf, 0xb5, 0xda, 0x0d, 0x7d, 0x7d, 0xdc, 0x57, 0xf3, 0xb8, 0x16, 0x37, 0xf8, 0x05,
	0x0f, 0x1b, 0x50, 0x3b, 0x43, 0x72, 0xba, 0x0e, 0x2b, 0xd9, 0x3b, 0xfb, 0x4b, 0xc2, 0x4b, 0x3e,
	0x4f, 0x01, 0x00, 0x00, 0xff, 0xff, 0xca, 0xb6, 0x71, 0x86, 0x20, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// InfoServerClient is the client API for InfoServer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type InfoServerClient interface {
	// Sends a greeting
	SetInfo(ctx context.Context, in *InfoRequest, opts ...grpc.CallOption) (*InfoReply, error)
	SetInfoStream(ctx context.Context, opts ...grpc.CallOption) (InfoServer_SetInfoStreamClient, error)
}

type infoServerClient struct {
	cc *grpc.ClientConn
}

func NewInfoServerClient(cc *grpc.ClientConn) InfoServerClient {
	return &infoServerClient{cc}
}

func (c *infoServerClient) SetInfo(ctx context.Context, in *InfoRequest, opts ...grpc.CallOption) (*InfoReply, error) {
	out := new(InfoReply)
	err := c.cc.Invoke(ctx, "/main.InfoServer/SetInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *infoServerClient) SetInfoStream(ctx context.Context, opts ...grpc.CallOption) (InfoServer_SetInfoStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &_InfoServer_serviceDesc.Streams[0], "/main.InfoServer/SetInfoStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &infoServerSetInfoStreamClient{stream}
	return x, nil
}

type InfoServer_SetInfoStreamClient interface {
	Send(*InfoRequest) error
	Recv() (*InfoReply, error)
	grpc.ClientStream
}

type infoServerSetInfoStreamClient struct {
	grpc.ClientStream
}

func (x *infoServerSetInfoStreamClient) Send(m *InfoRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *infoServerSetInfoStreamClient) Recv() (*InfoReply, error) {
	m := new(InfoReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// InfoServerServer is the server API for InfoServer service.
type InfoServerServer interface {
	// Sends a greeting
	SetInfo(context.Context, *InfoRequest) (*InfoReply, error)
	SetInfoStream(InfoServer_SetInfoStreamServer) error
}

func RegisterInfoServerServer(s *grpc.Server, srv InfoServerServer) {
	s.RegisterService(&_InfoServer_serviceDesc, srv)
}

func _InfoServer_SetInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InfoServerServer).SetInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/main.InfoServer/SetInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InfoServerServer).SetInfo(ctx, req.(*InfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _InfoServer_SetInfoStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(InfoServerServer).SetInfoStream(&infoServerSetInfoStreamServer{stream})
}

type InfoServer_SetInfoStreamServer interface {
	Send(*InfoReply) error
	Recv() (*InfoRequest, error)
	grpc.ServerStream
}

type infoServerSetInfoStreamServer struct {
	grpc.ServerStream
}

func (x *infoServerSetInfoStreamServer) Send(m *InfoReply) error {
	return x.ServerStream.SendMsg(m)
}

func (x *infoServerSetInfoStreamServer) Recv() (*InfoRequest, error) {
	m := new(InfoRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _InfoServer_serviceDesc = grpc.ServiceDesc{
	ServiceName: "main.InfoServer",
	HandlerType: (*InfoServerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SetInfo",
			Handler:    _InfoServer_SetInfo_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SetInfoStream",
			Handler:       _InfoServer_SetInfoStream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "info.proto",
}
