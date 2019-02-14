package obfs

import (
	"time"

	"github.com/rc452860/vnet/common/cache"
)

// golang dont support declare constant array
// so we use variable to replace it
var (
	DEFAULT_VERSION       = []byte{0x03, 0x03}
	DEFAULT_OVERHEAD      = 5
	DEFAULT_MAX_TIME_DIFF = 60 * 60 * 24
)

type ObfsAuthData struct {
	CLientData cache.Cache
	ClientId   []byte
	StartTIme  time.Time
	TIcketBuf  map[string][]byte
}

func NewObfsAuthData() *ObfsAuthData {
	return &ObfsAuthData{}
}

type ObfsTLS struct {
	plain
	HandshakeStatus int
	SendBuffer      []byte
	RecvBuffer      []byte
	ClientID        []byte
	MaxTimeDiff     int
	TLSVersion      []byte
	Overhead        int
}

func NewObfsTLS(method string) Plain {
	return &ObfsTLS{
		Method:          method,
		HandshakeStatus: 0,
		MaxTimeDiff:     DEFAULT_MAX_TIME_DIFF,
		TLSVersion:      DEFAULT_VERSION,
		Overhead:        DEFAULT_OVERHEAD,
	}
}

func (otls *ObfsTLS) InitData() []byte {
	panic("not implemented")
}

func (otls *ObfsTLS) GetOverhead(direction bool) int {
	panic("not implemented")
}

func (otls *ObfsTLS) GetServerInfo() ServerInfo {
	panic("not implemented")
}

func (otls *ObfsTLS) SetServerInfo(s ServerInfo) {
	panic("not implemented")
}

func (otls *ObfsTLS) ClientPreEncrypt(buf []byte) error {
	panic("not implemented")
}

func (otls *ObfsTLS) ClientEncode(buf []byte) error {
	panic("not implemented")
}

func (otls *ObfsTLS) ClientDecode(buf []byte) (bool, error) {
	panic("not implemented")
}

func (otls *ObfsTLS) ClientPostDecrypt(buf []byte) error {
	panic("not implemented")
}

func (otls *ObfsTLS) ServerPreEncrypt(buf []byte) error {
	panic("not implemented")
}

func (otls *ObfsTLS) ServerEncode(buf []byte) error {
	panic("not implemented")
}

func (otls *ObfsTLS) ServerDecode(buf []byte) (bool, bool, error) {
	panic("not implemented")
}

func (otls *ObfsTLS) ServerPostDecrypt(buf []byte) (bool, error) {
	panic("not implemented")
}

func (otls *ObfsTLS) ClientUDPPreEncrypt(buf []byte) error {
	panic("not implemented")
}

func (otls *ObfsTLS) ClientUDPPostDecrypt(buf []byte) error {
	panic("not implemented")
}

func (otls *ObfsTLS) ServerUDPPreEncrypt(buf []byte) error {
	panic("not implemented")
}

func (otls *ObfsTLS) ServerUDPPostDecrypt(buf []byte) (string, error) {
	panic("not implemented")
}

func (otls *ObfsTLS) Dispose() {
	panic("not implemented")
}

func (otls *ObfsTLS) GetHeadSize(buf []byte, defaultValue int) int {
	panic("not implemented")
}
