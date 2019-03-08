package obfs

type PlainFactory func(string) Plain

// Plain interface
type Plain interface {
	InitData() []byte
	GetMethod() string
	SetMethod(method string)
	GetOverhead(direction bool) int
	GetServerInfo() ServerInfo
	SetServerInfo(s ServerInfo)
	ClientPreEncrypt(buf []byte) ([]byte, error)
	ClientEncode(buf []byte) ([]byte, error)
	ClientDecode(buf []byte) ([]byte, bool, error)
	ClientPostDecrypt(buf []byte) ([]byte, error)
	ServerPreEncrypt(buf []byte) ([]byte, error)
	ServerEncode(buf []byte) ([]byte, error)
	ServerDecode(buf []byte) ([]byte, bool, bool, error)
	ServerPostDecrypt(buf []byte) ([]byte, bool, error)
	ClientUDPPreEncrypt(buf []byte) ([]byte, error)
	ClientUDPPostDecrypt(buf []byte) ([]byte, error)
	ServerUDPPreEncrypt(buf []byte) ([]byte, error)
	ServerUDPPostDecrypt(buf []byte) ([]byte, string, error)
	Dispose()
	GetHeadSize(buf []byte, defaultValue int) int
}

var (
	method_supported = make(map[string]PlainFactory)
)

func registerMethod(method string, factory PlainFactory) {
	method_supported[method] = factory
}
