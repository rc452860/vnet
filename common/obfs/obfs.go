package obfs

type PlainFactory func(string) Plain

// Plain interface
type Plain interface {
	InitData() []byte
	GetOverhead(direction bool) int
	GetServerInfo() ServerInfo
	SetServerInfo(s ServerInfo)
	ClientPreEncrypt(buf []byte) error
	ClientEncode(buf []byte) error
	ClientDecode(buf []byte) (bool, error)
	ClientPostDecrypt(buf []byte) error
	ServerPreEncrypt(buf []byte) error
	ServerEncode(buf []byte) error
	ServerDecode(buf []byte) (bool, bool, error)
	ServerPostDecrypt(buf []byte) (bool, error)
	ClientUDPPreEncrypt(buf []byte) error
	ClientUDPPostDecrypt(buf []byte) error
	ServerUDPPreEncrypt(buf []byte) error
	ServerUDPPostDecrypt(buf []byte) (string, error)
	Dispose()
	GetHeadSize(buf []byte, defaultValue int) int
}

var (
	method_supported = make(map[string]PlainFactory)
)

func registerMethod(method string, factory PlainFactory) {
	method_supported[method] = factory
}
