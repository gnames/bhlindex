// Code generated by protoc-gen-go. DO NOT EDIT.
// source: protob.proto

package protob

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type MatchType int32

const (
	MatchType_NONE            MatchType = 0
	MatchType_EXACT           MatchType = 1
	MatchType_CANONICAL_EXACT MatchType = 2
	MatchType_CANONICAL_FUZZY MatchType = 3
	MatchType_PARTIAL_EXACT   MatchType = 4
	MatchType_PARTIAL_FUZZY   MatchType = 5
)

var MatchType_name = map[int32]string{
	0: "NONE",
	1: "EXACT",
	2: "CANONICAL_EXACT",
	3: "CANONICAL_FUZZY",
	4: "PARTIAL_EXACT",
	5: "PARTIAL_FUZZY",
}
var MatchType_value = map[string]int32{
	"NONE":            0,
	"EXACT":           1,
	"CANONICAL_EXACT": 2,
	"CANONICAL_FUZZY": 3,
	"PARTIAL_EXACT":   4,
	"PARTIAL_FUZZY":   5,
}

func (x MatchType) String() string {
	return proto.EnumName(MatchType_name, int32(x))
}
func (MatchType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_protob_ce5f96dbb42af272, []int{0}
}

type Version struct {
	Value                string   `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Version) Reset()         { *m = Version{} }
func (m *Version) String() string { return proto.CompactTextString(m) }
func (*Version) ProtoMessage()    {}
func (*Version) Descriptor() ([]byte, []int) {
	return fileDescriptor_protob_ce5f96dbb42af272, []int{0}
}
func (m *Version) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Version.Unmarshal(m, b)
}
func (m *Version) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Version.Marshal(b, m, deterministic)
}
func (dst *Version) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Version.Merge(dst, src)
}
func (m *Version) XXX_Size() int {
	return xxx_messageInfo_Version.Size(m)
}
func (m *Version) XXX_DiscardUnknown() {
	xxx_messageInfo_Version.DiscardUnknown(m)
}

var xxx_messageInfo_Version proto.InternalMessageInfo

func (m *Version) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

type Void struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Void) Reset()         { *m = Void{} }
func (m *Void) String() string { return proto.CompactTextString(m) }
func (*Void) ProtoMessage()    {}
func (*Void) Descriptor() ([]byte, []int) {
	return fileDescriptor_protob_ce5f96dbb42af272, []int{1}
}
func (m *Void) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Void.Unmarshal(m, b)
}
func (m *Void) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Void.Marshal(b, m, deterministic)
}
func (dst *Void) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Void.Merge(dst, src)
}
func (m *Void) XXX_Size() int {
	return xxx_messageInfo_Void.Size(m)
}
func (m *Void) XXX_DiscardUnknown() {
	xxx_messageInfo_Void.DiscardUnknown(m)
}

var xxx_messageInfo_Void proto.InternalMessageInfo

type Title struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Path                 string   `protobuf:"bytes,2,opt,name=path,proto3" json:"path,omitempty"`
	Pages                []*Page  `protobuf:"bytes,3,rep,name=pages,proto3" json:"pages,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Title) Reset()         { *m = Title{} }
func (m *Title) String() string { return proto.CompactTextString(m) }
func (*Title) ProtoMessage()    {}
func (*Title) Descriptor() ([]byte, []int) {
	return fileDescriptor_protob_ce5f96dbb42af272, []int{2}
}
func (m *Title) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Title.Unmarshal(m, b)
}
func (m *Title) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Title.Marshal(b, m, deterministic)
}
func (dst *Title) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Title.Merge(dst, src)
}
func (m *Title) XXX_Size() int {
	return xxx_messageInfo_Title.Size(m)
}
func (m *Title) XXX_DiscardUnknown() {
	xxx_messageInfo_Title.DiscardUnknown(m)
}

var xxx_messageInfo_Title proto.InternalMessageInfo

func (m *Title) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Title) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func (m *Title) GetPages() []*Page {
	if m != nil {
		return m.Pages
	}
	return nil
}

type Page struct {
	Id                   string        `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Offset               int32         `protobuf:"varint,2,opt,name=offset,proto3" json:"offset,omitempty"`
	Names                []*NameString `protobuf:"bytes,3,rep,name=names,proto3" json:"names,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *Page) Reset()         { *m = Page{} }
func (m *Page) String() string { return proto.CompactTextString(m) }
func (*Page) ProtoMessage()    {}
func (*Page) Descriptor() ([]byte, []int) {
	return fileDescriptor_protob_ce5f96dbb42af272, []int{3}
}
func (m *Page) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Page.Unmarshal(m, b)
}
func (m *Page) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Page.Marshal(b, m, deterministic)
}
func (dst *Page) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Page.Merge(dst, src)
}
func (m *Page) XXX_Size() int {
	return xxx_messageInfo_Page.Size(m)
}
func (m *Page) XXX_DiscardUnknown() {
	xxx_messageInfo_Page.DiscardUnknown(m)
}

var xxx_messageInfo_Page proto.InternalMessageInfo

func (m *Page) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Page) GetOffset() int32 {
	if m != nil {
		return m.Offset
	}
	return 0
}

func (m *Page) GetNames() []*NameString {
	if m != nil {
		return m.Names
	}
	return nil
}

type NameString struct {
	Value                string    `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	Odds                 float32   `protobuf:"fixed32,2,opt,name=odds,proto3" json:"odds,omitempty"`
	Path                 string    `protobuf:"bytes,3,opt,name=path,proto3" json:"path,omitempty"`
	Curated              bool      `protobuf:"varint,4,opt,name=curated,proto3" json:"curated,omitempty"`
	EditDistance         int32     `protobuf:"varint,5,opt,name=edit_distance,json=editDistance,proto3" json:"edit_distance,omitempty"`
	EditDistanceStem     int32     `protobuf:"varint,6,opt,name=edit_distance_stem,json=editDistanceStem,proto3" json:"edit_distance_stem,omitempty"`
	Match                MatchType `protobuf:"varint,7,opt,name=match,proto3,enum=protob.MatchType" json:"match,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *NameString) Reset()         { *m = NameString{} }
func (m *NameString) String() string { return proto.CompactTextString(m) }
func (*NameString) ProtoMessage()    {}
func (*NameString) Descriptor() ([]byte, []int) {
	return fileDescriptor_protob_ce5f96dbb42af272, []int{4}
}
func (m *NameString) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NameString.Unmarshal(m, b)
}
func (m *NameString) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NameString.Marshal(b, m, deterministic)
}
func (dst *NameString) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NameString.Merge(dst, src)
}
func (m *NameString) XXX_Size() int {
	return xxx_messageInfo_NameString.Size(m)
}
func (m *NameString) XXX_DiscardUnknown() {
	xxx_messageInfo_NameString.DiscardUnknown(m)
}

var xxx_messageInfo_NameString proto.InternalMessageInfo

func (m *NameString) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func (m *NameString) GetOdds() float32 {
	if m != nil {
		return m.Odds
	}
	return 0
}

func (m *NameString) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func (m *NameString) GetCurated() bool {
	if m != nil {
		return m.Curated
	}
	return false
}

func (m *NameString) GetEditDistance() int32 {
	if m != nil {
		return m.EditDistance
	}
	return 0
}

func (m *NameString) GetEditDistanceStem() int32 {
	if m != nil {
		return m.EditDistanceStem
	}
	return 0
}

func (m *NameString) GetMatch() MatchType {
	if m != nil {
		return m.Match
	}
	return MatchType_NONE
}

func init() {
	proto.RegisterType((*Version)(nil), "protob.Version")
	proto.RegisterType((*Void)(nil), "protob.Void")
	proto.RegisterType((*Title)(nil), "protob.Title")
	proto.RegisterType((*Page)(nil), "protob.Page")
	proto.RegisterType((*NameString)(nil), "protob.NameString")
	proto.RegisterEnum("protob.MatchType", MatchType_name, MatchType_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// BHLIndexClient is the client API for BHLIndex service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type BHLIndexClient interface {
	Ver(ctx context.Context, in *Void, opts ...grpc.CallOption) (*Version, error)
	Titles(ctx context.Context, in *Void, opts ...grpc.CallOption) (BHLIndex_TitlesClient, error)
}

type bHLIndexClient struct {
	cc *grpc.ClientConn
}

func NewBHLIndexClient(cc *grpc.ClientConn) BHLIndexClient {
	return &bHLIndexClient{cc}
}

func (c *bHLIndexClient) Ver(ctx context.Context, in *Void, opts ...grpc.CallOption) (*Version, error) {
	out := new(Version)
	err := c.cc.Invoke(ctx, "/protob.BHLIndex/Ver", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bHLIndexClient) Titles(ctx context.Context, in *Void, opts ...grpc.CallOption) (BHLIndex_TitlesClient, error) {
	stream, err := c.cc.NewStream(ctx, &_BHLIndex_serviceDesc.Streams[0], "/protob.BHLIndex/Titles", opts...)
	if err != nil {
		return nil, err
	}
	x := &bHLIndexTitlesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type BHLIndex_TitlesClient interface {
	Recv() (*Title, error)
	grpc.ClientStream
}

type bHLIndexTitlesClient struct {
	grpc.ClientStream
}

func (x *bHLIndexTitlesClient) Recv() (*Title, error) {
	m := new(Title)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// BHLIndexServer is the server API for BHLIndex service.
type BHLIndexServer interface {
	Ver(context.Context, *Void) (*Version, error)
	Titles(*Void, BHLIndex_TitlesServer) error
}

func RegisterBHLIndexServer(s *grpc.Server, srv BHLIndexServer) {
	s.RegisterService(&_BHLIndex_serviceDesc, srv)
}

func _BHLIndex_Ver_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Void)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BHLIndexServer).Ver(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protob.BHLIndex/Ver",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BHLIndexServer).Ver(ctx, req.(*Void))
	}
	return interceptor(ctx, in, info, handler)
}

func _BHLIndex_Titles_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Void)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(BHLIndexServer).Titles(m, &bHLIndexTitlesServer{stream})
}

type BHLIndex_TitlesServer interface {
	Send(*Title) error
	grpc.ServerStream
}

type bHLIndexTitlesServer struct {
	grpc.ServerStream
}

func (x *bHLIndexTitlesServer) Send(m *Title) error {
	return x.ServerStream.SendMsg(m)
}

var _BHLIndex_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protob.BHLIndex",
	HandlerType: (*BHLIndexServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ver",
			Handler:    _BHLIndex_Ver_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Titles",
			Handler:       _BHLIndex_Titles_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "protob.proto",
}

func init() { proto.RegisterFile("protob.proto", fileDescriptor_protob_ce5f96dbb42af272) }

var fileDescriptor_protob_ce5f96dbb42af272 = []byte{
	// 406 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x92, 0x6d, 0x6b, 0x13, 0x41,
	0x10, 0xc7, 0x73, 0x0f, 0x7b, 0x49, 0xc6, 0xa4, 0x4d, 0x47, 0x91, 0xc5, 0x37, 0x86, 0x15, 0xf4,
	0x14, 0x29, 0x12, 0x3f, 0x41, 0x8c, 0x15, 0x03, 0xf5, 0x52, 0xb6, 0x31, 0xd4, 0x82, 0x84, 0x6d,
	0x76, 0x9b, 0x2e, 0xf4, 0x1e, 0xb8, 0xdb, 0x8a, 0x7e, 0x5b, 0x3f, 0x8a, 0xdc, 0x6e, 0xee, 0x4c,
	0xb5, 0xaf, 0x6e, 0xff, 0xbf, 0xf9, 0x31, 0xbb, 0xc3, 0x0d, 0x0c, 0x8a, 0x32, 0x37, 0xf9, 0xd5,
	0xb1, 0xfd, 0x60, 0xe4, 0x12, 0x7b, 0x0e, 0xdd, 0x95, 0x2a, 0x2b, 0x9d, 0x67, 0xf8, 0x04, 0xc8,
	0x0f, 0x71, 0x7b, 0xa7, 0xa8, 0x37, 0xf6, 0xe2, 0x3e, 0x77, 0x81, 0x45, 0x10, 0xae, 0x72, 0x2d,
	0xd9, 0x02, 0xc8, 0x52, 0x9b, 0x5b, 0x85, 0x07, 0xe0, 0x6b, 0xb9, 0x73, 0x7c, 0x2d, 0x11, 0x21,
	0x2c, 0x84, 0xb9, 0xa1, 0xbe, 0x25, 0xf6, 0x8c, 0x0c, 0x48, 0x21, 0xb6, 0xaa, 0xa2, 0xc1, 0x38,
	0x88, 0x1f, 0x4d, 0x06, 0xc7, 0xbb, 0xbb, 0xcf, 0xc4, 0x56, 0x71, 0x57, 0x62, 0x17, 0x10, 0xd6,
	0xf1, 0xbf, 0x7e, 0x4f, 0x21, 0xca, 0xaf, 0xaf, 0x2b, 0x65, 0x6c, 0x47, 0xc2, 0x77, 0x09, 0x63,
	0x20, 0x99, 0x48, 0xdb, 0x9e, 0xd8, 0xf4, 0x4c, 0x44, 0xaa, 0xce, 0x4d, 0xa9, 0xb3, 0x2d, 0x77,
	0x02, 0xfb, 0xed, 0x01, 0xfc, 0xa5, 0x0f, 0xcf, 0x55, 0x3f, 0x3b, 0x97, 0xb2, 0xb2, 0x97, 0xf8,
	0xdc, 0x9e, 0xdb, 0x51, 0x82, 0xbd, 0x51, 0x28, 0x74, 0x37, 0x77, 0xa5, 0x30, 0x4a, 0xd2, 0x70,
	0xec, 0xc5, 0x3d, 0xde, 0x44, 0x7c, 0x01, 0x43, 0x25, 0xb5, 0x59, 0x4b, 0x5d, 0x19, 0x91, 0x6d,
	0x14, 0x25, 0xf6, 0xbd, 0x83, 0x1a, 0x7e, 0xdc, 0x31, 0x7c, 0x0b, 0x78, 0x4f, 0x5a, 0x57, 0x46,
	0xa5, 0x34, 0xb2, 0xe6, 0x68, 0xdf, 0x3c, 0x37, 0x2a, 0xc5, 0x57, 0x40, 0x52, 0x61, 0x36, 0x37,
	0xb4, 0x3b, 0xf6, 0xe2, 0x83, 0xc9, 0x51, 0x33, 0xe3, 0x97, 0x1a, 0x2e, 0x7f, 0x15, 0x8a, 0xbb,
	0xfa, 0x9b, 0x02, 0xfa, 0x2d, 0xc3, 0x1e, 0x84, 0xc9, 0x22, 0x39, 0x19, 0x75, 0xb0, 0x0f, 0xe4,
	0xe4, 0x62, 0x3a, 0x5b, 0x8e, 0x3c, 0x7c, 0x0c, 0x87, 0xb3, 0x69, 0xb2, 0x48, 0xe6, 0xb3, 0xe9,
	0xe9, 0xda, 0x41, 0xff, 0x3e, 0xfc, 0xf4, 0xf5, 0xf2, 0xf2, 0xdb, 0x28, 0xc0, 0x23, 0x18, 0x9e,
	0x4d, 0xf9, 0x72, 0xde, 0x7a, 0xe1, 0x3e, 0x72, 0x16, 0x99, 0x7c, 0x87, 0xde, 0x87, 0xcf, 0xa7,
	0xf3, 0x4c, 0xaa, 0x9f, 0xf8, 0x12, 0x82, 0x95, 0x2a, 0xb1, 0xfd, 0xad, 0xf5, 0x82, 0x3c, 0x3b,
	0x6c, 0x93, 0xdb, 0x27, 0xd6, 0xc1, 0xd7, 0x10, 0xd9, 0x9d, 0xa9, 0xfe, 0x51, 0x87, 0x4d, 0xb2,
	0x55, 0xd6, 0x79, 0xe7, 0x5d, 0xb9, 0x7d, 0x7c, 0xff, 0x27, 0x00, 0x00, 0xff, 0xff, 0x02, 0x0f,
	0x43, 0xb1, 0xa6, 0x02, 0x00, 0x00,
}