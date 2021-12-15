// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: github.com/bittorrent/go-btfs/protos/shard/shard.proto

package shardpb

import (
	fmt "fmt"
	golang_proto "github.com/golang/protobuf/proto"
	guard "github.com/tron-us/go-btfs-common/protos/guard"
	_ "github.com/tron-us/protobuf/gogoproto"
	proto "github.com/tron-us/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = golang_proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type Status struct {
	Status               string   `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty" pg:"status"`
	Message              string   `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty" pg:"message"`
	XXX_NoUnkeyedLiteral struct{} `json:"-" pg:"-"`
	XXX_unrecognized     []byte   `json:"-" pg:"-"`
	XXX_sizecache        int32    `json:"-" pg:"-"`
}

func (m *Status) Reset()         { *m = Status{} }
func (m *Status) String() string { return proto.CompactTextString(m) }
func (*Status) ProtoMessage()    {}
func (*Status) Descriptor() ([]byte, []int) {
	return fileDescriptor_4f69c1466e2dc4b4, []int{0}
}
func (m *Status) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Status) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Status.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Status) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Status.Merge(m, src)
}
func (m *Status) XXX_Size() int {
	return m.Size()
}
func (m *Status) XXX_DiscardUnknown() {
	xxx_messageInfo_Status.DiscardUnknown(m)
}

var xxx_messageInfo_Status proto.InternalMessageInfo

func (m *Status) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *Status) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (*Status) XXX_MessageName() string {
	return "shard.Status"
}

type AdditionalInfo struct {
	Info                 string   `protobuf:"bytes,1,opt,name=info,proto3" json:"info,omitempty" pg:"info"`
	XXX_NoUnkeyedLiteral struct{} `json:"-" pg:"-"`
	XXX_unrecognized     []byte   `json:"-" pg:"-"`
	XXX_sizecache        int32    `json:"-" pg:"-"`
}

func (m *AdditionalInfo) Reset()         { *m = AdditionalInfo{} }
func (m *AdditionalInfo) String() string { return proto.CompactTextString(m) }
func (*AdditionalInfo) ProtoMessage()    {}
func (*AdditionalInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_4f69c1466e2dc4b4, []int{1}
}
func (m *AdditionalInfo) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *AdditionalInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_AdditionalInfo.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *AdditionalInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AdditionalInfo.Merge(m, src)
}
func (m *AdditionalInfo) XXX_Size() int {
	return m.Size()
}
func (m *AdditionalInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_AdditionalInfo.DiscardUnknown(m)
}

var xxx_messageInfo_AdditionalInfo proto.InternalMessageInfo

func (m *AdditionalInfo) GetInfo() string {
	if m != nil {
		return m.Info
	}
	return ""
}

func (*AdditionalInfo) XXX_MessageName() string {
	return "shard.AdditionalInfo"
}

type SignedContracts struct {
	SignedEscrowContract []byte          `protobuf:"bytes,1,opt,name=signed_escrow_contract,json=signedEscrowContract,proto3" json:"signed_escrow_contract,omitempty" pg:"signed_escrow_contract"`
	SignedGuardContract  *guard.Contract `protobuf:"bytes,2,opt,name=signed_guard_contract,json=signedGuardContract,proto3" json:"signed_guard_contract,omitempty" pg:"signed_guard_contract"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-" pg:"-"`
	XXX_unrecognized     []byte          `json:"-" pg:"-"`
	XXX_sizecache        int32           `json:"-" pg:"-"`
}

func (m *SignedContracts) Reset()         { *m = SignedContracts{} }
func (m *SignedContracts) String() string { return proto.CompactTextString(m) }
func (*SignedContracts) ProtoMessage()    {}
func (*SignedContracts) Descriptor() ([]byte, []int) {
	return fileDescriptor_4f69c1466e2dc4b4, []int{2}
}
func (m *SignedContracts) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SignedContracts) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SignedContracts.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SignedContracts) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SignedContracts.Merge(m, src)
}
func (m *SignedContracts) XXX_Size() int {
	return m.Size()
}
func (m *SignedContracts) XXX_DiscardUnknown() {
	xxx_messageInfo_SignedContracts.DiscardUnknown(m)
}

var xxx_messageInfo_SignedContracts proto.InternalMessageInfo

func (m *SignedContracts) GetSignedEscrowContract() []byte {
	if m != nil {
		return m.SignedEscrowContract
	}
	return nil
}

func (m *SignedContracts) GetSignedGuardContract() *guard.Contract {
	if m != nil {
		return m.SignedGuardContract
	}
	return nil
}

func (*SignedContracts) XXX_MessageName() string {
	return "shard.SignedContracts"
}
func init() {
	proto.RegisterType((*Status)(nil), "shard.Status")
	golang_proto.RegisterType((*Status)(nil), "shard.Status")
	proto.RegisterType((*AdditionalInfo)(nil), "shard.AdditionalInfo")
	golang_proto.RegisterType((*AdditionalInfo)(nil), "shard.AdditionalInfo")
	proto.RegisterType((*SignedContracts)(nil), "shard.SignedContracts")
	golang_proto.RegisterType((*SignedContracts)(nil), "shard.SignedContracts")
}

func init() {
	proto.RegisterFile("github.com/bittorrent/go-btfs/protos/shard/shard.proto", fileDescriptor_4f69c1466e2dc4b4)
}
func init() {
	golang_proto.RegisterFile("github.com/bittorrent/go-btfs/protos/shard/shard.proto", fileDescriptor_4f69c1466e2dc4b4)
}

var fileDescriptor_4f69c1466e2dc4b4 = []byte{
	// 321 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x51, 0xbd, 0x4e, 0xf3, 0x30,
	0x14, 0x95, 0x3f, 0x7d, 0xb4, 0xc2, 0xfc, 0x54, 0x32, 0x50, 0x55, 0x1d, 0x2c, 0x54, 0x31, 0xb0,
	0x34, 0x41, 0x94, 0xa9, 0x1b, 0xad, 0x10, 0x62, 0xab, 0x92, 0x8d, 0xa5, 0xca, 0xaf, 0x6b, 0x89,
	0xf8, 0x56, 0xb1, 0x23, 0x5e, 0x82, 0x07, 0xe0, 0x71, 0x18, 0x3b, 0xf2, 0x08, 0x90, 0xbc, 0x04,
	0x23, 0xf2, 0x75, 0x22, 0x16, 0x58, 0x6e, 0xee, 0x39, 0xe7, 0x9e, 0xa3, 0xe4, 0x84, 0xce, 0x84,
	0x34, 0x9b, 0x2a, 0xf6, 0x12, 0x28, 0x7c, 0x53, 0x82, 0x9a, 0x56, 0xda, 0x17, 0x30, 0x8d, 0x4d,
	0xae, 0xfd, 0x6d, 0x09, 0x06, 0xb4, 0xaf, 0x37, 0x51, 0x99, 0xba, 0xe9, 0x21, 0xc5, 0xf6, 0x10,
	0x8c, 0xe7, 0x7f, 0x7b, 0xa7, 0x09, 0x14, 0x05, 0xa8, 0x2e, 0x42, 0x54, 0x36, 0x02, 0xa7, 0x8b,
	0x18, 0x5f, 0xfd, 0xe2, 0x45, 0x25, 0xae, 0x72, 0x5f, 0x80, 0x00, 0x04, 0xb8, 0x39, 0xc7, 0x64,
	0x4e, 0x7b, 0xa1, 0x89, 0x4c, 0xa5, 0xd9, 0x90, 0xf6, 0x34, 0x6e, 0x23, 0x72, 0x4e, 0x2e, 0xf7,
	0x83, 0x16, 0xb1, 0x11, 0xed, 0x17, 0x99, 0xd6, 0x91, 0xc8, 0x46, 0xff, 0x50, 0xe8, 0xe0, 0xe4,
	0x82, 0x1e, 0xdf, 0xa6, 0xa9, 0x34, 0x12, 0x54, 0xf4, 0xf4, 0xa0, 0x72, 0x60, 0x8c, 0xfe, 0x97,
	0x2a, 0x87, 0x36, 0x01, 0xf7, 0xc9, 0x0b, 0xa1, 0x83, 0x50, 0x0a, 0x95, 0xa5, 0x4b, 0x50, 0xa6,
	0x8c, 0x12, 0xa3, 0xd9, 0x0d, 0x1d, 0x6a, 0xa4, 0xd6, 0x99, 0x4e, 0x4a, 0x78, 0x5e, 0x27, 0xad,
	0x84, 0xce, 0xc3, 0xe0, 0xd4, 0xa9, 0x77, 0x28, 0x76, 0x36, 0xb6, 0xa4, 0x67, 0xad, 0x0b, 0xbf,
	0xf9, 0xc7, 0x64, 0xdf, 0xeb, 0xe0, 0x7a, 0xe0, 0xb9, 0x2a, 0xba, 0xfb, 0xe0, 0xc4, 0x5d, 0xdf,
	0x5b, 0xb6, 0x23, 0x17, 0x8b, 0xaf, 0x4f, 0x4e, 0x76, 0x35, 0x27, 0xef, 0x35, 0x27, 0x1f, 0x35,
	0x27, 0xaf, 0x0d, 0x27, 0x6f, 0x0d, 0x27, 0xbb, 0x86, 0x13, 0x7a, 0x24, 0xc1, 0xb3, 0x4d, 0x7b,
	0xf8, 0x2f, 0x16, 0x34, 0xb4, 0x8f, 0x95, 0x6d, 0x6a, 0x45, 0x1e, 0xfb, 0x48, 0x6e, 0xe3, 0xb8,
	0x87, 0xdd, 0xcd, 0xbe, 0x03, 0x00, 0x00, 0xff, 0xff, 0xc3, 0x92, 0x64, 0xa3, 0xe7, 0x01, 0x00,
	0x00,
}

func (m *Status) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Status) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Status) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Message) > 0 {
		i -= len(m.Message)
		copy(dAtA[i:], m.Message)
		i = encodeVarintShard(dAtA, i, uint64(len(m.Message)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Status) > 0 {
		i -= len(m.Status)
		copy(dAtA[i:], m.Status)
		i = encodeVarintShard(dAtA, i, uint64(len(m.Status)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *AdditionalInfo) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AdditionalInfo) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *AdditionalInfo) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Info) > 0 {
		i -= len(m.Info)
		copy(dAtA[i:], m.Info)
		i = encodeVarintShard(dAtA, i, uint64(len(m.Info)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *SignedContracts) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SignedContracts) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SignedContracts) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.SignedGuardContract != nil {
		{
			size, err := m.SignedGuardContract.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintShard(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.SignedEscrowContract) > 0 {
		i -= len(m.SignedEscrowContract)
		copy(dAtA[i:], m.SignedEscrowContract)
		i = encodeVarintShard(dAtA, i, uint64(len(m.SignedEscrowContract)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintShard(dAtA []byte, offset int, v uint64) int {
	offset -= sovShard(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func NewPopulatedStatus(r randyShard, easy bool) *Status {
	this := &Status{}
	this.Status = string(randStringShard(r))
	this.Message = string(randStringShard(r))
	if !easy && r.Intn(10) != 0 {
		this.XXX_unrecognized = randUnrecognizedShard(r, 3)
	}
	return this
}

func NewPopulatedAdditionalInfo(r randyShard, easy bool) *AdditionalInfo {
	this := &AdditionalInfo{}
	this.Info = string(randStringShard(r))
	if !easy && r.Intn(10) != 0 {
		this.XXX_unrecognized = randUnrecognizedShard(r, 2)
	}
	return this
}

func NewPopulatedSignedContracts(r randyShard, easy bool) *SignedContracts {
	this := &SignedContracts{}
	v1 := r.Intn(100)
	this.SignedEscrowContract = make([]byte, v1)
	for i := 0; i < v1; i++ {
		this.SignedEscrowContract[i] = byte(r.Intn(256))
	}
	if r.Intn(5) != 0 {
		this.SignedGuardContract = guard.NewPopulatedContract(r, easy)
	}
	if !easy && r.Intn(10) != 0 {
		this.XXX_unrecognized = randUnrecognizedShard(r, 3)
	}
	return this
}

type randyShard interface {
	Float32() float32
	Float64() float64
	Int63() int64
	Int31() int32
	Uint32() uint32
	Intn(n int) int
}

func randUTF8RuneShard(r randyShard) rune {
	ru := r.Intn(62)
	if ru < 10 {
		return rune(ru + 48)
	} else if ru < 36 {
		return rune(ru + 55)
	}
	return rune(ru + 61)
}
func randStringShard(r randyShard) string {
	v2 := r.Intn(100)
	tmps := make([]rune, v2)
	for i := 0; i < v2; i++ {
		tmps[i] = randUTF8RuneShard(r)
	}
	return string(tmps)
}
func randUnrecognizedShard(r randyShard, maxFieldNumber int) (dAtA []byte) {
	l := r.Intn(5)
	for i := 0; i < l; i++ {
		wire := r.Intn(4)
		if wire == 3 {
			wire = 5
		}
		fieldNumber := maxFieldNumber + r.Intn(100)
		dAtA = randFieldShard(dAtA, r, fieldNumber, wire)
	}
	return dAtA
}
func randFieldShard(dAtA []byte, r randyShard, fieldNumber int, wire int) []byte {
	key := uint32(fieldNumber)<<3 | uint32(wire)
	switch wire {
	case 0:
		dAtA = encodeVarintPopulateShard(dAtA, uint64(key))
		v3 := r.Int63()
		if r.Intn(2) == 0 {
			v3 *= -1
		}
		dAtA = encodeVarintPopulateShard(dAtA, uint64(v3))
	case 1:
		dAtA = encodeVarintPopulateShard(dAtA, uint64(key))
		dAtA = append(dAtA, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))
	case 2:
		dAtA = encodeVarintPopulateShard(dAtA, uint64(key))
		ll := r.Intn(100)
		dAtA = encodeVarintPopulateShard(dAtA, uint64(ll))
		for j := 0; j < ll; j++ {
			dAtA = append(dAtA, byte(r.Intn(256)))
		}
	default:
		dAtA = encodeVarintPopulateShard(dAtA, uint64(key))
		dAtA = append(dAtA, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))
	}
	return dAtA
}
func encodeVarintPopulateShard(dAtA []byte, v uint64) []byte {
	for v >= 1<<7 {
		dAtA = append(dAtA, uint8(uint64(v)&0x7f|0x80))
		v >>= 7
	}
	dAtA = append(dAtA, uint8(v))
	return dAtA
}
func (m *Status) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Status)
	if l > 0 {
		n += 1 + l + sovShard(uint64(l))
	}
	l = len(m.Message)
	if l > 0 {
		n += 1 + l + sovShard(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *AdditionalInfo) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Info)
	if l > 0 {
		n += 1 + l + sovShard(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *SignedContracts) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.SignedEscrowContract)
	if l > 0 {
		n += 1 + l + sovShard(uint64(l))
	}
	if m.SignedGuardContract != nil {
		l = m.SignedGuardContract.Size()
		n += 1 + l + sovShard(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovShard(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozShard(x uint64) (n int) {
	return sovShard(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Status) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowShard
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Status: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Status: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Status", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowShard
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthShard
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthShard
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Status = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Message", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowShard
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthShard
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthShard
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Message = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipShard(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthShard
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthShard
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *AdditionalInfo) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowShard
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: AdditionalInfo: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AdditionalInfo: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Info", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowShard
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthShard
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthShard
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Info = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipShard(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthShard
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthShard
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *SignedContracts) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowShard
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: SignedContracts: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SignedContracts: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SignedEscrowContract", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowShard
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthShard
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthShard
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SignedEscrowContract = append(m.SignedEscrowContract[:0], dAtA[iNdEx:postIndex]...)
			if m.SignedEscrowContract == nil {
				m.SignedEscrowContract = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SignedGuardContract", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowShard
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthShard
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthShard
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.SignedGuardContract == nil {
				m.SignedGuardContract = &guard.Contract{}
			}
			if err := m.SignedGuardContract.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipShard(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthShard
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthShard
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipShard(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowShard
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowShard
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowShard
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthShard
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupShard
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthShard
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthShard        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowShard          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupShard = fmt.Errorf("proto: unexpected end of group")
)
