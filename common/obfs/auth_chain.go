package obfs

import (
	"bytes"

	"github.com/rc452860/vnet/utils/binaryx"
)

const (
	MAX_INT  = (1 << 64) - 1
	MOV_MASK = (1 << (64 - 23)) - 1
)

type XorShift128Plus struct {
	V0 uint64
	V1 uint64
}

func NewXorShift128Plus() *XorShift128Plus {
	return &XorShift128Plus{
		V0: 0,
		V1: 0,
	}
}

func (xs1p *XorShift128Plus) Next() uint64 {
	x := xs1p.V0
	y := xs1p.V1
	xs1p.V0 = y
	x ^= ((x & MOV_MASK) << 23)
	x ^= (y ^ (x >> 17) ^ (y >> 26))
	xs1p.V1 = x
	return (x + y) & MAX_INT
}

func (xs1p *XorShift128Plus) InitFromBin(bin []byte) {
	if len(bin) < 16 {
		bin = conbineToBytes(bin, bytes.Repeat([]byte{byte(0x00)}, 16))
	}
	xs1p.V0 = binaryx.LEBytesToUint64(bin[0:8])
	xs1p.V1 = binaryx.LEBytesToUint64(bin[8:16])
}

func (xs1p *XorShift128Plus) InitFromBinLen(bin []byte, length int) {
	if len(bin) < 16 {
		bin = conbineToBytes(bin, bytes.Repeat([]byte{byte(0x00)}, 16))
	}
	xs1p.V0 = binaryx.LEBytesToUint64(conbineToBytes(binaryx.LEUInt16ToBytes(uint16(length)), bin[2:8]))
	xs1p.V1 = binaryx.LEBytesToUint64(bin[8:16])
	for i := 0; i < 4; i++ {
		xs1p.Next()
	}
}

/*--------------------------AuthBase----------------------------*/

type AuthBase struct {
	Plain
	Method             string
	NoCompatibleMethod string
	Overhead           int
	RawTrans           bool
}

func NewAuthBase(method string) *AuthBase {
	return &AuthBase{
		Plain:    NewPlain(method),
		Method:   method,
		Overhead: 4,
	}
}

func (authBase *AuthBase) GetOverhead(direction bool) int {
	return authBase.Overhead
}

func (authBase *AuthBase) NotMatchReturn(buf []byte) ([]byte, bool) {
	authBase.RawTrans = true
	authBase.Overhead = 0
	if authBase.GetMethod() == authBase.NoCompatibleMethod {
		return bytes.Repeat([]byte{byte('E')}, 2048), false
	}
	return buf, false
}

type AuthChainA struct {
	AuthBase
	RecvBuf        []byte
	UnintLen       int
	HasSentHeader  bool
	HasRecvHeader  bool
	ClientId       int
	ConnectionId   int
	MaxTimeDif     int
	Salt           []byte
	PackId         int
	RecvId         int
	UserId         int
	UserIdNum      int
	UserKey        []byte
	ClientOverhead int
	LastClientHash []byte
	LastServerHash []byte
	RandomClient   XorShift128Plus
	RandomServer   XorShift128Plus
	
}
