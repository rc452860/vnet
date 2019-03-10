package obfs

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
)

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

func hmacmd5(key, data []byte) []byte {
	mac := hmac.New(md5.New, key)
	mac.Write(data)
	return mac.Sum(nil)
}
