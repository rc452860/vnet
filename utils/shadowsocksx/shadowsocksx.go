package shadowsocksx

import (
	"crypto/sha1"
	"golang.org/x/crypto/hkdf"
	"io"
)

func HKDF_SHA1(secret, salt, info, key []byte) error {
	_, err := io.ReadFull(hkdf.New(sha1.New, secret, salt, info), key)
	return err
}
