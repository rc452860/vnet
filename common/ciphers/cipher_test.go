package ciphers

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func BenchmarkDecrypto(t *testing.B) {
	for i := 0; i < t.N; i++ {
		key, _ := hex.DecodeString("b36d331451a61eb2d76860e00c347396")
		iv, _ := hex.DecodeString("3b272b460e99f22a314d5f7335c00e6e")
		cipher, _ := NewCipher(OP_DECRYPT, "aes-128-gcm", key, iv)
		data, _ := hex.DecodeString("b1c5a0097220c3db3e0b23b98c8b3cddd292fd9ddf769e8b7494de52f6d958b19712f10f08a096878b")
		cipher.Decrypt(data)
	}
}

func ExampleDecrypto() {
	key, _ := hex.DecodeString("b36d331451a61eb2d76860e00c347396")
	iv, _ := hex.DecodeString("3b272b460e99f22a314d5f7335c00e6e")
	cipher, _ := NewCipher(OP_DECRYPT, "aes-128-gcm", key, iv)
	data, _ := hex.DecodeString("b1c5a0097220c3db3e0b23b98c8b3cddd292fd9ddf769e8b7494de52f6d958b19712f10f08a096878b")
	result, err := cipher.Decrypt(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(hex.EncodeToString(result))
	//Output:
	//6675636b796f75
}

func ExampleEncrypto() {
	key, _ := hex.DecodeString("b36d331451a61eb2d76860e00c347396")
	iv, _ := hex.DecodeString("3b272b460e99f22a314d5f7335c00e6e")
	encipher, _ := NewCipher(OP_ENCRYPT, "aes-128-gcm", key, iv)
	result, err := encipher.Encrypt([]byte("fuckyou"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(hex.EncodeToString(result))

	decipher, _ := NewCipher(OP_DECRYPT, "aes-128-gcm", key, iv)
	result, err = decipher.Decrypt(result)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(hex.EncodeToString(result))
	//Output:
	//b1c5a0097220c3db3e0b23b98c8b3cddd292fd9ddf769e8b7494de52f6d958b19712f10f08a096878b
	//6675636b796f75
}
