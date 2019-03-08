package ciphers

type ICipher interface {
	GetKey() []byte
	GetIv() []byte
}
