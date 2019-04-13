package block

import (
	"crypto/aes"
	"crypto/cipher"
)

func init() {
	registerBlockCiphers("aes-128-cbc", &aes_cbc{16, 16})
}

type aes_cbc struct {
	keyLen int
	ivLen  int
}

func (a *aes_cbc) KeyLen() int {
	return a.keyLen
}
func (a *aes_cbc) IVLen() int {
	return a.ivLen
}
func (a *aes_cbc) NewBlock(key, iv []byte,encryptor bool) (cipher.BlockMode, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if encryptor{
		return cipher.NewCBCEncrypter(block, iv),nil
	}else{
		return cipher.NewCBCDecrypter(block, iv),nil
	}
}