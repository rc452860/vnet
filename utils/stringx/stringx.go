package stringx

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/sirupsen/logrus"
	"strconv"
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

func MustUnquote(s string) string {
	r, err := strconv.Unquote(s)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":    err,
			"string": s,
		}).Error("MustUnquote error: ", err)
		return ""
	}
	return r
}

// Convert string like \u4f60\u597d to utf-8 encode
// \u4f60\u597d means 你好(hello)
func UnicodeToUtf8(s string) string {
	slen := len(s)
	i := 0
	stringBuffer := new(bytes.Buffer)
	for i < slen {
		if s[i] == 92 && (s[i+1] == 85 || s[i+1] == 117) {
			temp,err:=strconv.ParseInt(s[i+2:i+6],16,32)
			if err != nil{
				panic(err)
			}
			stringBuffer.WriteString(fmt.Sprintf("%c", temp))
			i += 6
			continue
		}else{
			stringBuffer.WriteByte(s[i])
			i++
			continue
		}
	}
	return stringBuffer.String()
}

// Convert string like \u4f60\u597d to utf-8 encode
// \u4f60\u597d means 你好(hello)
func BUnicodeToUtf8(s []byte) string {
	slen := len(s)
	i := 0
	stringBuffer := new(bytes.Buffer)
	for i < slen {
		if s[i] == 92 && (s[i+1] == 85 || s[i+1] == 117) {
			temp,err:=strconv.ParseInt(string(s[i+2:i+6]),16,32)
			if err != nil{
				panic(err)
			}
			stringBuffer.WriteString(fmt.Sprintf("%c", temp))
			i += 6
			continue
		}else{
			stringBuffer.WriteByte(s[i])
			i++
			continue
		}
	}
	return stringBuffer.String()
}

