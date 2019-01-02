package utils

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

type SS_Scheme struct {
	Method   string `json:"method"`
	Password string `json:"password"`
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
}

func Parse_SIP002_URI_Scheme(s string) (*SS_Scheme, error) {
	if !strings.HasPrefix(s, "ss://") {
		return nil, errors.New("not sip002 scheme!")
	}
	sInfo := s[5:]
	decodeByte, err := base64.StdEncoding.DecodeString(sInfo)
	if err != nil {
		return nil, err
	}
	decodeStr := string(decodeByte)
	fmt.Println(decodeStr)
	return nil, err
}
