package obfs

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/rc452860/vnet/common/cache"
	"github.com/rc452860/vnet/common/log"
	"github.com/rc452860/vnet/utils/arrayx"
	"github.com/rc452860/vnet/utils/randomx"
	"github.com/rc452860/vnet/utils/stringx"
)

// golang dont support declare constant array
// so we use variable to replace it
var (
	DEFAULT_VERSION       = []byte{0x03, 0x03}
	DEFAULT_OVERHEAD      = 5
	DEFAULT_MAX_TIME_DIFF = 60 * 60 * 24
)

type ObfsAuthData struct {
	ServerInfo
	ClientData *cache.Cache
	ClientID   []byte
	StartTIme  int
	TicketBuf  map[string][]byte
}

func NewObfsAuthData() *ObfsAuthData {
	return &ObfsAuthData{
		ServerInfo: NewServerInfo(),
		TicketBuf:  make(map[string][]byte),
		ClientID:   randomx.RandomBytes(32),
		ClientData: cache.New(60 * 5 * time.Second),
	}
}

type ObfsTLS struct {
	Plain
	*ObfsAuthData
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
		Plain:           NewPlain(method),
		HandshakeStatus: 0,
		MaxTimeDiff:     DEFAULT_MAX_TIME_DIFF,
		TLSVersion:      DEFAULT_VERSION,
		Overhead:        DEFAULT_OVERHEAD,
		ObfsAuthData:    NewObfsAuthData(),
	}
}

func (otls *ObfsTLS) InitData() []byte {
	panic("not implemented")
}

func (otls *ObfsTLS) GetOverhead(direction bool) int {
	return otls.Overhead
}

func (otls *ObfsTLS) GetServerInfo() ServerInfo {
	return otls.Plain.GetServerInfo()
}

func (otls *ObfsTLS) SetServerInfo(s ServerInfo) {
	otls.Plain.SetServerInfo(s)
}

func (otls *ObfsTLS) ClientPreEncrypt(buf []byte) ([]byte, error) {
	return otls.Plain.ClientPreEncrypt(buf)
}

func (otls *ObfsTLS) ClientEncode(buf []byte) ([]byte, error) {
	if otls.HandshakeStatus == -1 {
		return buf, nil
	}

	if (otls.HandshakeStatus & 8) == 8 {
		ret := []byte{}
		for len(buf) > 2048 {
			left := float64(binary.BigEndian.Uint16(randomx.RandomBytes(2))%4096 + 100)
			right := float64(len(buf))
			size := uint16(math.Min(left, right))
			ret = conbineToBytes(ret, byte(0x17), otls.TLSVersion, uint16(size), buf[:size])
			buf = buf[size:]
		}
		if len(buf) > 0 {
			ret = conbineToBytes(ret, byte(0x17), otls.TLSVersion, uint16(len(buf)), buf)
		}
		return ret, nil
	}
	if len(buf) > 0 {
		otls.SendBuffer = conbineToBytes(otls.SendBuffer, byte(0x17), otls.TLSVersion, uint16(len(buf)), buf)
	}
	if otls.HandshakeStatus == 0 {
		otls.HandshakeStatus = 1
		data := new(bytes.Buffer)
		ext := new(bytes.Buffer)
		binary.Write(data, binary.BigEndian, otls.TLSVersion)
		binary.Write(data, binary.BigEndian, otls.packAuthData(otls.ObfsAuthData.ClientID))
		binary.Write(data, binary.BigEndian, byte(0x20))
		binary.Write(data, binary.BigEndian, otls.ObfsAuthData.ClientID)
		binary.Write(data, binary.BigEndian, []byte{0x0, 0x1c, 0xc0, 0x2b, 0xc0, 0x2f, 0xcc, 0xa9, 0xcc, 0xa8, 0xcc, 0x14, 0xcc, 0x13, 0xc0, 0xa,
			0xc0, 0x14, 0xc0, 0x9, 0xc0, 0x13, 0x0, 0x9c, 0x0, 0x35, 0x0, 0x2f, 0x0, 0xa})
		binary.Write(data, binary.BigEndian, []byte{0x0, 0x01})

		binary.Write(ext, binary.BigEndian, []byte{0xff, 0x1, 0x0, 0x1, 0x0})

		var host string
		if otls.GetServerInfo().GetObfsParam() != "" {
			host = otls.GetServerInfo().GetObfsParam()
		} else {
			host = otls.GetServerInfo().GetHost()
		}

		if host != "" && stringx.IsDigit(string(host[len(host)-1])) {
			host = ""
		}

		hosts := strings.Split(host, ",")
		host = randomx.RandomStringsChoice(hosts)

		binary.Write(ext, binary.BigEndian, otls.sni(host))
		binary.Write(ext, binary.BigEndian, []byte{0x00, 0x17, 0x00, 0x00})

		if otls.ObfsAuthData.TicketBuf[host] == nil {
			otls.ObfsAuthData.TicketBuf[host] = randomx.RandomBytes((int(randomx.Uint16())%17 + 8) * 16)
		}
		binary.Write(ext, binary.BigEndian, conbineToBytes(
			[]byte{0x00, 0x23},
			len(otls.ObfsAuthData.TicketBuf[host]),
			otls.ObfsAuthData.TicketBuf[host]))

		binary.Write(ext, binary.BigEndian, MustHexDecode("000d001600140601060305010503040104030301030302010203"))
		binary.Write(ext, binary.BigEndian, MustHexDecode("000500050100000000"))
		binary.Write(ext, binary.BigEndian, MustHexDecode("00120000"))
		binary.Write(ext, binary.BigEndian, MustHexDecode("75500000"))
		binary.Write(ext, binary.BigEndian, MustHexDecode("000b00020100"))
		binary.Write(ext, binary.BigEndian, MustHexDecode("000a0006000400170018"))

		binary.Write(data, binary.BigEndian, conbineToBytes(
			uint16(ext.Len()),
			ext.Bytes()))

		result := conbineToBytes([]byte{0x01, 0x00}, uint16(data.Len()), data.Bytes())
		result = conbineToBytes([]byte{0x16, 0x03, 0x01}, uint16(len(result)), result)
		return result, nil
	} else if otls.HandshakeStatus == 1 && len(buf) == 0 {
		data := conbineToBytes(byte(0x14), otls.TLSVersion, []byte{0x00, 0x01, 0x01}) //ChangeCipherSpec
		data = conbineToBytes(data, byte(0x16), otls.TLSVersion, []byte{0x00, 0x20}, randomx.RandomBytes(22))
		data = conbineToBytes(data, hmacsha1(conbineToBytes(otls.GetServerInfo().GetKey(), otls.ObfsAuthData.ClientID), data)[:10])
		ret := conbineToBytes(data, otls.SendBuffer)
		otls.SendBuffer = []byte{}
		otls.HandshakeStatus = 8
		return ret, nil
	}

	return []byte{}, nil
}

//ClientDecode buffer_to_recv, is_need_to_encode_and_send_back
func (otls *ObfsTLS) ClientDecode(buf []byte) ([]byte, bool, error) {
	if otls.HandshakeStatus == -1 {
		return buf, false, nil
	}
	if otls.HandshakeStatus == 8 {
		ret := new(bytes.Buffer)
		otls.RecvBuffer = conbineToBytes(otls.RecvBuffer, buf)
		for len(otls.RecvBuffer) > 5 {
			if int(otls.RecvBuffer[0]) != 0x17 {
				log.Error("data = %s", hex.EncodeToString(otls.RecvBuffer))
				return nil, false, errors.New("server_decode appdata error")
			}
			size := binary.BigEndian.Uint16(otls.RecvBuffer[3:5])
			if len(otls.RecvBuffer) < int(size)+5 {
				break
			}
			buf = otls.RecvBuffer[5 : size+5]
			binary.Write(ret, binary.BigEndian, buf)
			otls.RecvBuffer = otls.RecvBuffer[size+5:]
		}
		return ret.Bytes(), false, nil
	}

	if len(buf) < 11+32+1+32 {
		return nil, false, errors.New("client_decode data error")
	}

	verify := buf[11:33]
	if !bytes.Equal(hmacsha1(conbineToBytes(otls.GetServerInfo().GetKey(), otls.ObfsAuthData.ClientID), verify)[:10], buf[33:43]) {
		return nil, false, errors.New("client_decode data error")
	}
	if !bytes.Equal(hmacsha1(conbineToBytes(otls.GetServerInfo().GetKey(), otls.ObfsAuthData.ClientID), buf[:len(buf)-10])[:10], buf[len(buf)-10:]) {
		return nil, false, errors.New("client_decode data error")
	}
	return []byte{}, true, nil
}

func (otls *ObfsTLS) ClientPostDecrypt(buf []byte) ([]byte, error) {
	panic("not implemented")
}

func (otls *ObfsTLS) ServerPreEncrypt(buf []byte) ([]byte, error) {
	panic("not implemented")
}

func (otls *ObfsTLS) ServerEncode(buf []byte) ([]byte, error) {
	if otls.HandshakeStatus == -1 {
		return buf, nil
	}
	if (otls.HandshakeStatus & 8) == 8 {
		ret := new(bytes.Buffer)
		for len(buf) > 2048 {
			size := uint16(math.Min(float64(randomx.Uint16()%4096+100), float64(len(buf))))
			binary.Write(ret, binary.BigEndian, conbineToBytes(
				byte(0x17),
				otls.TLSVersion,
				size,
				buf[:size]))
			buf = buf[size:]
		}
		if len(buf) > 0 {
			binary.Write(ret, binary.BigEndian, conbineToBytes(
				byte(0x17),
				otls.TLSVersion,
				uint16(len(buf)),
				buf))
			return ret.Bytes(), nil
		}
	}

	otls.HandshakeStatus |= 8
	data := conbineToBytes(otls.TLSVersion, otls.packAuthData(otls.ClientID), byte(0x20), otls.ClientID, MustHexDecode("c02f000005ff01000100"))
	data = conbineToBytes([]byte{0x20, 0x00}, uint16(len(data)), data) // server hello
	data = conbineToBytes(byte(0x16), otls.TLSVersion, uint16(len(data)), data)
	if int(randomx.Float64Range(0, 8)) < 1 {
		ticket := randomx.RandomBytes(int((randomx.Uint16()%164)*2) + 64)
		ticket = conbineToBytes(uint16(len(ticket)+4), []byte{0x04, 0x00}, len(ticket), ticket)
		data = conbineToBytes(data, byte(0x16), otls.TLSVersion, ticket) // New session ticket
	}
	data = conbineToBytes(data, byte(0x14), otls.TLSVersion, []byte{0x00, 0x01, 0x01}) // ChangeCipherSpec
	finishLen := randomx.RandomIntChoice([]int{32, 40})
	data = conbineToBytes(data, byte(0x16), otls.TLSVersion, uint16(finishLen), randomx.RandomBytes(finishLen-10))
	data = conbineToBytes(data, hmacsha1(conbineToBytes(otls.GetServerInfo().GetKey(), otls.ClientID), data)[:10])
	if buf != nil && len(buf) != 0 {
		tmp, err := otls.ServerEncode(buf)
		if err != nil {
			return nil, err
		}
		data = conbineToBytes(data, tmp)
	}
	return data, nil
}

// ServerDecode return buffer_to_recv, is_need_decrypt, is_need_to_encode_and_send_back
func (otls *ObfsTLS) ServerDecode(buf []byte) ([]byte, bool, bool, error) {
	if otls.HandshakeStatus == -1 {
		return buf, true, false, nil
	}
	if (otls.HandshakeStatus & 4) == 4 {
		ret := new(bytes.Buffer)
		otls.RecvBuffer = conbineToBytes(otls.RecvBuffer, buf)
		for len(otls.RecvBuffer) > 5 {
			if int(otls.RecvBuffer[0]) != 0x17 || int(otls.RecvBuffer[1]) != 0x03 ||
				int(otls.RecvBuffer[2]) != 0x03 {
				log.Error("data = %s", hex.EncodeToString(otls.RecvBuffer))
				return nil, false, false, errors.New("server_decode appdata error")
			}

			size := binary.BigEndian.Uint16(otls.RecvBuffer[3:5])
			if len(otls.RecvBuffer) < int(size)+5 {
				break
			}
			binary.Write(ret, binary.BigEndian, otls.RecvBuffer[5:size+5])
			otls.RecvBuffer = otls.RecvBuffer[size+5:]
		}
		return ret.Bytes(), true, false, nil
	}

	if (otls.HandshakeStatus & 1) == 1 {
		otls.RecvBuffer = conbineToBytes(otls.RecvBuffer, buf)
		buf = otls.RecvBuffer
		verify := buf
		if len(buf) < 11 {
			return nil, false, false, errors.New("server_decode data error")
		}

		if !matchBegin(buf, conbineToBytes(byte(0x14), otls.TLSVersion, []byte{0x00, 0x01, 0x01})) {
			return nil, false, false, errors.New("server_decode data error")
		}
		buf = buf[6:]
		if !matchBegin(buf, conbineToBytes(byte(0x16), otls.TLSVersion, byte(0x00))) {
			return nil, false, false, errors.New("server_decode data error")
		}

		verifyLen := binary.BigEndian.Uint16(buf[3:5]) + 1 // 11 - 10
		if len(verify) < int(verifyLen)+10 {
			return []byte{}, false, false, nil
		}
		if !bytes.Equal(hmacsha1(conbineToBytes(otls.GetServerInfo().GetKey(), otls.ClientID), verify[:verifyLen])[:10], verify[verifyLen:verifyLen+10]) {
			return nil, false, false, errors.New("server_decode data error")
		}
		otls.RecvBuffer = verify[verifyLen+10:]
		// status := otls.HandshakeStatus
		otls.HandshakeStatus |= 4
		return otls.ServerDecode([]byte{})
	}
	otls.RecvBuffer = conbineToBytes(otls.RecvBuffer, buf)
	buf = otls.RecvBuffer
	originBuf := buf
	if len(buf) < 3 {
		return []byte{}, false, false, nil
	}
	if !matchBegin(buf, []byte{0x16, 0x03, 0x01}) {
		return otls.DecodeErrorReturn(originBuf)
	}
	buf = buf[3:]
	headerLen := binary.BigEndian.Uint16(buf[:2])
	if headerLen > uint16(len(buf))-2 {
		return []byte{}, false, false, nil
	}

	otls.RecvBuffer = otls.RecvBuffer[headerLen+5:]
	otls.HandshakeStatus = 1
	buf = buf[2 : headerLen+2]
	if !matchBegin(buf, []byte{0x01, 0x00}) {
		log.Info("tls_auth not client hello message")
		return otls.DecodeErrorReturn(originBuf)
	}
	buf = buf[2:]
	if binary.BigEndian.Uint16(buf) != uint16(len(buf))-2 {
		log.Info("tls_auth wrong message size")
		return otls.DecodeErrorReturn(originBuf)
	}
	buf = buf[2:]
	if !matchBegin(buf, otls.TLSVersion) {
		log.Info("tls_auth wrong tls version")
		return otls.DecodeErrorReturn(originBuf)
	}
	buf = buf[2:]
	verifyId := buf[:32]
	buf = buf[32:]
	sessionLen := int8(buf[0])
	if sessionLen < 32 {
		log.Info("tls_auth wrong sessionid_len")
		return otls.DecodeErrorReturn(originBuf)
	}
	sessionId := buf[1 : sessionLen+1]
	buf = buf[sessionLen+1:]
	otls.ClientID = sessionId
	sha1 := hmacsha1(conbineToBytes(otls.GetServerInfo().GetKey(), sessionId), verifyId[:22])[:10]
	utcTime := int(binary.BigEndian.Uint32(verifyId[:4]))
	timeDif := int(uint32(time.Now().Unix()) - uint32(utcTime))

	if otls.GetServerInfo().GetObfsParam() != "" {
		dif, err := strconv.Atoi(otls.GetServerInfo().GetObfsParam())
		if err == nil {
			otls.MaxTimeDiff = dif
		}
	}

	if otls.MaxTimeDiff > 0 &&
		(timeDif < -otls.MaxTimeDiff ||
			timeDif > otls.MaxTimeDiff || int32(utcTime-otls.ObfsAuthData.StartTIme) < int32(otls.MaxTimeDiff/2)) {
		log.Info("tls_auth wrong time")
		return otls.DecodeErrorReturn(originBuf)
	}

	if !bytes.Equal(sha1, verifyId[22:]) {
		log.Info("tls_auth wrong sha1")
		return otls.DecodeErrorReturn(originBuf)
	}

	if otls.ClientData.Get(string(verifyId[:22])) != nil {
		log.Info("replay attack detect, id = %s", hex.EncodeToString(verifyId))
		return otls.DecodeErrorReturn(originBuf)
	}

	otls.ClientData.Put(string(verifyId[:22]), sessionId, time.Duration(DEFAULT_MAX_TIME_DIFF)*time.Second)
	if len(otls.RecvBuffer) >= 11 {
		ret, _, _, _ := otls.ServerDecode([]byte{})
		return ret, true, true, nil
	}
	return []byte{}, false, true, nil
}

func (otls *ObfsTLS) ServerPostDecrypt(buf []byte) ([]byte, bool, error) {
	panic("not implemented")
}

func (otls *ObfsTLS) ClientUDPPreEncrypt(buf []byte) ([]byte, error) {
	panic("not implemented")
}

func (otls *ObfsTLS) ClientUDPPostDecrypt(buf []byte) ([]byte, error) {
	panic("not implemented")
}

func (otls *ObfsTLS) ServerUDPPreEncrypt(buf []byte) ([]byte, error) {
	panic("not implemented")
}

func (otls *ObfsTLS) ServerUDPPostDecrypt(buf []byte) ([]byte, string, error) {
	panic("not implemented")
}

func (otls *ObfsTLS) Dispose() {
	panic("not implemented")
}

func (otls *ObfsTLS) GetHeadSize(buf []byte, defaultValue int) int {
	panic("not implemented")
}

func (otls *ObfsTLS) packAuthData(clientId []byte) []byte {
	dataBuf := new(bytes.Buffer)
	binary.Write(dataBuf, binary.BigEndian, uint32(time.Now().Unix()&0xFFFFFFFF))
	binary.Write(dataBuf, binary.BigEndian, randomx.RandomBytes(18))
	binary.Write(dataBuf, binary.BigEndian, hmacsha1(conbineToBytes(otls.Plain.GetServerInfo().GetKey(), clientId), dataBuf.Bytes())[:10])
	return dataBuf.Bytes()
}

func (otls *ObfsTLS) sni(host string) []byte {
	url := []byte(host)
	data := conbineToBytes([]byte{0x00}, uint16(len(url)), url)
	data = conbineToBytes([]byte{0x00, 0x00}, uint16(len(data)+2), uint16(len(data)), data)
	return data
}

func (otls *ObfsTLS) DecodeErrorReturn(buf []byte) ([]byte, bool, bool, error) {
	otls.HandshakeStatus = -1
	if otls.Overhead > 0 {
		otls.GetServerInfo().SetOverhead(otls.GetServerInfo().GetOverhead() - otls.Overhead)
	}
	otls.Overhead = 0
	if arrayx.FindStringInArray(otls.Plain.GetMethod(), []string{"tls1.2_ticket_auth", "tls1.2_ticket_fastauth"}) {
		return bytes.Repeat([]byte{byte('E')}, 2048), false, false, nil
	}

	return buf, true, false, nil
}

func conbineToBytes(data ...interface{}) []byte {
	buf := new(bytes.Buffer)
	for _, item := range data {
		binary.Write(buf, binary.BigEndian, item)
	}
	return buf.Bytes()
}

func MustHexDecode(data string) []byte {
	result, err := hex.DecodeString(data)
	if err != nil {
		return []byte{}
	}
	return result
}

func hmacsha1(key, data []byte) []byte {
	mac := hmac.New(sha1.New, key)
	mac.Write(data)
	return mac.Sum(nil)
}

func matchBegin(str1, str2 []byte) bool {
	if len(str1) >= len(str2) {
		if bytes.Equal(str1[:len(str2)], str2) {
			return true
		}
	}
	return false
}
