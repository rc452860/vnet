package stream

import (
	"crypto/cipher"
	"crypto/rc4"
)

func init() {
	registerStreamCiphers("rc4", &rc4_cryptor{16, 0})

}

type rc4_cryptor struct {
	keyLen int
	ivLen  int
}

func (a *rc4_cryptor) KeyLen() int {
	return a.keyLen
}
func (a *rc4_cryptor) IVLen() int {
	return a.ivLen
}

func (a *rc4_cryptor) NewStream(key, iv []byte, _ int) (cipher.Stream, error) {
	block, err := rc4.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return block, nil
}
