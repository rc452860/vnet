package stream

import (
	"crypto/cipher"

	"golang.org/x/crypto/blowfish"
)

func init() {
	registerStreamCiphers("bf-cfb", &bf_cfb{16, 8})
}

type bf_cfb struct {
	keyLen int
	ivLen  int
}

func (a *bf_cfb) KeyLen() int {
	return a.keyLen
}
func (a *bf_cfb) IVLen() int {
	return a.ivLen
}
func (a *bf_cfb) NewStream(key, iv []byte) (cipher.Stream, error) {
	block, err := blowfish.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewCFBEncrypter(block, iv), nil
}
