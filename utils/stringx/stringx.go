package stringx

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"strings"
)

// IsDigit judement data is any of 1234567890
func IsDigit(data string) bool {
	if len(data) != 1 {
		return false
	}
	if strings.IndexAny(data, "1234567890") != -1 {
		return true
	}
	return false
}

func U2S(form string) (to string, err error) {
	bs, err := hex.DecodeString(strings.Replace(form, `\u`, ``, -1))
	if err != nil {
		return
	}
	for i, bl, br, r := 0, len(bs), bytes.NewReader(bs), uint16(0); i < bl; i += 2 {
		binary.Read(br, binary.BigEndian, &r)
		to += string(r)
	}
	return
}

func MustU2S(from string) string {
	result, err := U2S(from)
	if err != nil {
		return ""
	}
	return result
}
