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
	// ServerDecode return buffer_to_recv, is_need_decrypt, is_need_to_encode_and_send_back
	ServerDecode(buf []byte) ([]byte, bool, bool, error)
	ServerPostDecrypt(buf []byte) ([]byte, bool, error)
	ClientUDPPreEncrypt(buf []byte) ([]byte, error)
	ClientUDPPostDecrypt(buf []byte) ([]byte, error)
	ServerUDPPreEncrypt(buf,uid []byte) ([]byte, error)
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

func GetObfs(method string) Plain{
	return method_supported[method](method)
}
