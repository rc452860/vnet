package aead

import (
	"crypto/cipher"

	"golang.org/x/crypto/chacha20poly1305"
)

func init() {
	registerAEADCiphers("chacha20-ietf-poly1305", &chacha20IetfPoly1305{32, 32, 12, 16})
}

type chacha20IetfPoly1305 struct {
	keySize   int
	saltSize  int
	nonceSize int
	tagSize   int
}

func (c *chacha20IetfPoly1305) KeySize() int {
	return c.keySize
}

func (c *chacha20IetfPoly1305) SaltSize() int {
	return c.saltSize
}

func (c *chacha20IetfPoly1305) NonceSize() int {
	return c.nonceSize
}

func (c *chacha20IetfPoly1305) NewAEAD(key []byte, salt []byte) (cipher.AEAD, error) {
	subkey := make([]byte, c.KeySize())
	HKDF_SHA1(key, salt, []byte("ss-subkey"), subkey)
	return chacha20poly1305.New(subkey)
}
