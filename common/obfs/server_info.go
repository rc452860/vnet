package obfs

import "net"

const (
	NETWORK_MTU      = 1500
	TCP_MSS          = 1460
	BUF_SIZE         = 32 * 1024
	UDP_MAX_BUF_SIZE = 65536
	DEFAULT_HEAD_LEN = 30
)

type ServerInfo interface {
	GetHost() string
	SetHost(host string)
	GetPort() int
	SetPort(port int)
	GetClient() net.IP
	SetClient(client net.IP)
	GetProtocolParam() string
	SetProtocolParam(protocolParam string)
	GetObfsParam() string
	SetObfsParam(obfsParam string)
	SetIv(iv []byte)
	GetIv() []byte
	SetRecvIv(iv []byte)
	GetRecvIv() []byte
	SetKeyStr(key string)
	GetKeyStr() string
	SetHeadLen(len int)
	GetHeadLen() int
	SetTCPMss(mss int)
	GetTCPMss() int
	SetBufferSize(size int)
	GetBufferSize() int
	SetOverhead(size int)
	GetOverhead() int
}

type serverInfo struct {
	Host          string
	Port          int
	Client        net.IP
	ClientPort    int
	ProtocolParam string
	ObfsParam     string
	Iv            []byte
	RecvIv        []byte
	KeyStr        string
	Key           []byte
	HeadLen       int
	TCPMss        int
	BufferSize    int
	Overhead      int
}

// InitServerInfo init ServerInfo default value
func NewServerInfo() ServerInfo {
	return &serverInfo{
		TCPMss:  TCP_MSS,
		HeadLen: DEFAULT_HEAD_LEN,
	}
}

func (s *serverInfo) GetHost() string {
	panic("not implemented")
}

func (s *serverInfo) SetHost(host string) {
	panic("not implemented")
}

func (s *serverInfo) GetPort() int {
	panic("not implemented")
}

func (s *serverInfo) SetPort(port int) {
	panic("not implemented")
}

func (s *serverInfo) GetClient() net.IP {
	panic("not implemented")
}

func (s *serverInfo) SetClient(client net.IP) {
	panic("not implemented")
}

func (s *serverInfo) GetProtocolParam() string {
	panic("not implemented")
}

func (s *serverInfo) SetProtocolParam(protocolParam string) {
	panic("not implemented")
}

func (s *serverInfo) GetObfsParam() string {
	panic("not implemented")
}

func (s *serverInfo) SetObfsParam(obfsParam string) {
	panic("not implemented")
}

func (s *serverInfo) SetIv(iv []byte) {
	panic("not implemented")
}

func (s *serverInfo) GetIv() []byte {
	panic("not implemented")
}

func (s *serverInfo) SetRecvIv(iv []byte) {
	panic("not implemented")
}

func (s *serverInfo) GetRecvIv() []byte {
	panic("not implemented")
}

func (s *serverInfo) SetKeyStr(key string) {
	panic("not implemented")
}

func (s *serverInfo) GetKeyStr() string {
	panic("not implemented")
}

func (s *serverInfo) SetHeadLen(len int) {
	panic("not implemented")
}

func (s *serverInfo) GetHeadLen() int {
	panic("not implemented")
}

func (s *serverInfo) SetTCPMss(mss int) {
	panic("not implemented")
}

func (s *serverInfo) GetTCPMss() int {
	panic("not implemented")
}

func (s *serverInfo) SetBufferSize(size int) {
	panic("not implemented")
}

func (s *serverInfo) GetBufferSize() int {
	panic("not implemented")
}

func (s *serverInfo) SetOverhead(size int) {
	panic("not implemented")
}

func (s *serverInfo) GetOverhead() int {
	panic("not implemented")
}
