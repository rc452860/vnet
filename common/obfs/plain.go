package obfs

func init() {
	registerMethod("plain", NewPlain)
	registerMethod("origin", NewPlain)
}

type plain struct {
	ServerInfo
	Method string
}

// NewPlain construct a new plain and initliza default value
func NewPlain(method string) Plain {
	return &plain{
		Method: method,
	}
}

func (p *plain) GetMethod() string {
	return p.Method
}

func (p *plain) SetMethod(method string) {
	p.Method = method
}

func (p *plain) InitData() []byte {
	return []byte{}
}

func (p *plain) GetOverhead(direction bool) int {
	return 0
}

func (p *plain) GetServerInfo() ServerInfo {
	return p.ServerInfo
}

func (p *plain) SetServerInfo(s ServerInfo) {
	p.ServerInfo = s
}

func (p *plain) ClientPreEncrypt(buf []byte) ([]byte, error) {
	return buf, nil
}

func (p *plain) ClientEncode(buf []byte) ([]byte, error) {
	return buf, nil
}

//ClientDecode buffer_to_recv, is_need_to_encode_and_send_back
func (p *plain) ClientDecode(buf []byte) ([]byte, bool, error) {
	return buf, false, nil
}

func (p *plain) ClientPostDecrypt(buf []byte) ([]byte, error) {
	return buf, nil
}

func (p *plain) ServerPreEncrypt(buf []byte) ([]byte, error) {
	return buf, nil
}

func (p *plain) ServerEncode(buf []byte) ([]byte, error) {
	return buf, nil
}

//buffer_to_recv, is_need_decrypt, is_need_to_encode_and_send_back
func (p *plain) ServerDecode(buf []byte) ([]byte, bool, bool, error) {
	return buf, true, false, nil
}

func (p *plain) ServerPostDecrypt(buf []byte) ([]byte, bool, error) {
	return buf, false, nil
}

func (p *plain) ClientUDPPreEncrypt(buf []byte) ([]byte, error) {
	return buf, nil
}

func (p *plain) ClientUDPPostDecrypt(buf []byte) ([]byte, error) {
	return buf, nil
}

func (p *plain) ServerUDPPreEncrypt(buf,uid []byte) ([]byte, error) {
	return buf, nil
}

func (p *plain) ServerUDPPostDecrypt(buf []byte) ([]byte, string, error) {
	return buf, "", nil
}

func (p *plain) Dispose() {

}

func (p *plain) GetHeadSize(buf []byte, defaultValue int) int {
	if len(buf) < 2 {
		return defaultValue
	}
	headType := int(buf[0]) & 0x7
	switch headType {
	case 1:
		return 7
	case 4:
		return 19
	case 3:
		return 4 + int(buf[1])
	}
	return defaultValue
}
