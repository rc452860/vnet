package block

import "crypto/cipher"

type IBlockCipher interface {
	KeyLen() int
	IVLen() int
	NewBlock(key []byte, iv []byte,encryptor bool) (cipher.BlockMode, error)
}

var blockCiphers = make(map[string]IBlockCipher)

func registerBlockCiphers(method string, c IBlockCipher) {
	blockCiphers[method] = c
}

func GetBlockCiphers() map[string]IBlockCipher {
	return blockCiphers
}

func GetBlockCipher(method string) IBlockCipher {
	return blockCiphers[method]
}

