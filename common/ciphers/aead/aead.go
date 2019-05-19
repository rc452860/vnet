package aead

import "crypto/cipher"

type IAEADCipher interface {
	KeySize() int
	SaltSize() int
	NonceSize() int
	NewAEAD(key []byte, salt []byte, decryptOrEncrypt int) (cipher.AEAD, error)
}

var aeadCiphers = make(map[string]IAEADCipher)

func registerAEADCiphers(method string, c IAEADCipher) {
	aeadCiphers[method] = c
}

func GetAEADCiphers() map[string]IAEADCipher {
	return aeadCiphers
}

func GetAEADCipher(method string) IAEADCipher {
	return aeadCiphers[method]
}
