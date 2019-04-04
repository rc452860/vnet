package stream

import (
	"crypto/cipher"
	"crypto/des"
)

func init() {
	registerStreamCiphers("des-cfb", &des_cfb{8, 8})
}

type des_cfb struct {
	keyLen int
	ivLen  int
}

func (a *des_cfb) KeyLen() int {
	return a.keyLen
}
func (a *des_cfb) IVLen() int {
	return a.ivLen
}
func (a *des_cfb) NewStream(key, iv []byte) (cipher.Stream, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewCFBEncrypter(block, iv), nil
}
