package stream

import (
	"crypto/cipher"

	"golang.org/x/crypto/cast5"
)

func init() {
	registerStreamCiphers("cast5-cfb", &cast5_cfb{16, 8})
}

type cast5_cfb struct {
	keyLen int
	ivLen  int
}

func (a *cast5_cfb) KeyLen() int {
	return a.keyLen
}
func (a *cast5_cfb) IVLen() int {
	return a.ivLen
}
func (a *cast5_cfb) NewStream(key, iv []byte, decryptOrEncrypt int) (cipher.Stream, error) {
	block, err := cast5.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if decryptOrEncrypt == 0 {
		return cipher.NewCFBEncrypter(block, iv), nil
	} else {
		return cipher.NewCFBDecrypter(block, iv), nil
	}
}
