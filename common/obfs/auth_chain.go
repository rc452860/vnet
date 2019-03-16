package obfs

import (
	"bytes"

	"github.com/rc452860/vnet/common/cache"
	"github.com/rc452860/vnet/common/obfs"
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

type ObfsAUthChainData struct {
	Name         stirng
	UserId       cache.Cache
	LastClientId []byte
	ConnectionId int
	MaxClient    int
	MaxBuffer    int
}

/*----------------------------------AuthChainA----------------------------------*/

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
	RandomClient   *XorShift128Plus
	RandomServer   *XorShift128Plus
}

func NewAuthChainA() *AuthChainA {
	return &AuthChainA{
		AuthBase: {
			RawTrans:           false,
			Overhead:           4,
			NoCompatibleMethod: "auth_chain_a",
		},
		RecvBuf:        []byte{},
		UnintLen:       2800,
		HasRecvHeader:  false,
		HasSentHeader:  false,
		ClientId:       0,
		ConnectionId:   0,
		MaxTimeDif:     60 * 60 * 24,
		Salt:           []byte("auth_chain_a"),
		PackId:         1,
		RecvId:         1,
		UserIdNum:      0,
		ClientOverhead: 4,
		LastClientHash: []byte{},
		LastServerHash: []byte{},
		RandomClient:   NewXorShift128Plus(),
		RandomServer:   NewXorShift128Plus(),
	}
}

func (a *AuthChainA) InitData() []byte {
	panic("not implemented")
}

func (a *AuthChainA) GetMethod() string {
	panic("not implemented")
}

func (a *AuthChainA) SetMethod(method string) {
	panic("not implemented")
}

func (a *AuthChainA) GetOverhead(direction bool) int {
	panic("not implemented")
}

func (a *AuthChainA) GetServerInfo() obfs.ServerInfo {
	panic("not implemented")
}

func (a *AuthChainA) SetServerInfo(s obfs.ServerInfo) {
	panic("not implemented")
}

func (a *AuthChainA) ClientPreEncrypt(buf []byte) ([]byte, error) {
	panic("not implemented")
}

func (a *AuthChainA) ClientEncode(buf []byte) ([]byte, error) {
	panic("not implemented")
}

func (a *AuthChainA) ClientDecode(buf []byte) ([]byte, bool, error) {
	panic("not implemented")
}

func (a *AuthChainA) ClientPostDecrypt(buf []byte) ([]byte, error) {
	panic("not implemented")
}

func (a *AuthChainA) ServerPreEncrypt(buf []byte) ([]byte, error) {
	panic("not implemented")
}

func (a *AuthChainA) ServerEncode(buf []byte) ([]byte, error) {
	panic("not implemented")
}

func (a *AuthChainA) ServerDecode(buf []byte) ([]byte, bool, bool, error) {
	panic("not implemented")
}

func (a *AuthChainA) ServerPostDecrypt(buf []byte) ([]byte, bool, error) {
	panic("not implemented")
}

func (a *AuthChainA) ClientUDPPreEncrypt(buf []byte) ([]byte, error) {
	panic("not implemented")
}

func (a *AuthChainA) ClientUDPPostDecrypt(buf []byte) ([]byte, error) {
	panic("not implemented")
}

func (a *AuthChainA) ServerUDPPreEncrypt(buf []byte) ([]byte, error) {
	panic("not implemented")
}

func (a *AuthChainA) ServerUDPPostDecrypt(buf []byte) ([]byte, string, error) {
	panic("not implemented")
}

func (a *AuthChainA) Dispose() {
	panic("not implemented")
}

func (a *AuthChainA) GetHeadSize(buf []byte, defaultValue int) int {
	panic("not implemented")
}
