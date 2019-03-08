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
	SetKey(key []byte)
	GetKey() []byte
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
	return s.Host
}

func (s *serverInfo) SetHost(host string) {
	s.Host = host
}

func (s *serverInfo) GetPort() int {
	return s.Port
}

func (s *serverInfo) SetPort(port int) {
	s.Port = port
}

func (s *serverInfo) GetClient() net.IP {
	return s.Client
}

func (s *serverInfo) SetClient(client net.IP) {
	s.Client = client
}

func (s *serverInfo) GetProtocolParam() string {
	return s.ProtocolParam
}

func (s *serverInfo) SetProtocolParam(protocolParam string) {
	s.ProtocolParam = protocolParam
}

func (s *serverInfo) GetObfsParam() string {
	return s.ObfsParam
}

func (s *serverInfo) SetObfsParam(obfsParam string) {
	s.ObfsParam = obfsParam
}

func (s *serverInfo) SetIv(iv []byte) {
	s.Iv = iv
}

func (s *serverInfo) GetIv() []byte {
	return s.Iv
}

func (s *serverInfo) SetRecvIv(iv []byte) {
	s.RecvIv = iv
}

func (s *serverInfo) GetRecvIv() []byte {
	return s.RecvIv
}

func (s *serverInfo) SetKeyStr(key string) {
	s.KeyStr = key
}

func (s *serverInfo) GetKeyStr() string {
	return s.KeyStr
}

func (s *serverInfo) SetKey(key []byte) {
	s.Key = key
}

func (s *serverInfo) GetKey() []byte {
	return s.Key
}

func (s *serverInfo) SetHeadLen(len int) {
	s.HeadLen = len
}

func (s *serverInfo) GetHeadLen() int {
	return s.HeadLen
}

func (s *serverInfo) SetTCPMss(mss int) {
	s.TCPMss = mss
}

func (s *serverInfo) GetTCPMss() int {
	return s.TCPMss
}

func (s *serverInfo) SetBufferSize(size int) {
	s.BufferSize = size
}

func (s *serverInfo) GetBufferSize() int {
	return s.BufferSize
}

func (s *serverInfo) SetOverhead(size int) {
	s.Overhead = size
}

func (s *serverInfo) GetOverhead() int {
	return s.Overhead
}
