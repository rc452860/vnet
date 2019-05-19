package stream

import "crypto/cipher"

type IStreamCipher interface {
	KeyLen() int
	IVLen() int
	// decryptOrEncrypt 0: encrypt 1: decrypt
	NewStream(key []byte, iv []byte, decryptOrEncrypt int) (cipher.Stream, error)
}

var streamCiphers = make(map[string]IStreamCipher)

func registerStreamCiphers(method string, c IStreamCipher) {
	streamCiphers[method] = c
}

func GetStreamCiphers() map[string]IStreamCipher {
	return streamCiphers
}

func GetStreamCipher(method string) IStreamCipher {
	return streamCiphers[method]
}
