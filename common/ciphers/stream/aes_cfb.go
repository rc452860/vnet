package stream

import (
	"crypto/aes"
	"crypto/cipher"
)

func init() {
	registerStreamCiphers("aes-128-cfb", &aes_cfb{16, 16})
	registerStreamCiphers("aes-192-cfb", &aes_cfb{24, 16})
	registerStreamCiphers("aes-256-cfb", &aes_cfb{32, 16})
}

type aes_cfb struct {
	keyLen int
	ivLen  int
}

func (a *aes_cfb) KeyLen() int {
	return a.keyLen
}
func (a *aes_cfb) IVLen() int {
	return a.ivLen
}
func (a *aes_cfb) NewStream(key, iv []byte, decryptOrEncrypt int) (cipher.Stream, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if decryptOrEncrypt == 0 {
		return cipher.NewCFBEncrypter(block, iv), nil
	} else {
		return cipher.NewCFBDecrypter(block, iv), nil
	}
}
