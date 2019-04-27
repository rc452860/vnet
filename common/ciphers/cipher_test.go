package ciphers

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"github.com/rc452860/vnet/utils/bytesx"
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

func TestEncryptorEncryptForStream(t *testing.T) {
	encryptor, err := NewEncryptor("aes-128-cfb", "killer")
	if err != nil {
		fmt.Println(err)
		return
	}
	raw := []byte("abc")
	fmt.Println(hex.EncodeToString(raw))
	result, err := encryptor.Encrypt(raw)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("ciphertext result: " + hex.EncodeToString(result))
	result, err = encryptor.Decrypt(result)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("plain text: " + hex.EncodeToString(result))
	if hex.EncodeToString(result) != "616263" {
		t.Fail()
	}
	//Output:
}

func TestEncryptorEncryptForAead(t *testing.T) {
	encryptor, err := NewEncryptor("aes-128-gcm", "killer")
	if err != nil {
		fmt.Println(err)
		return
	}
	raw := []byte("abc")
	fmt.Println(hex.EncodeToString(raw))
	result, err := encryptor.Encrypt(raw)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("ciphertext result: " + hex.EncodeToString(result))

	result, err = encryptor.Decrypt(result)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("plain text: " + hex.EncodeToString(result))
	if hex.EncodeToString(result) != "616263" {
		t.Fail()
	}
}

func TestEncryptor_Decrypt_CBC(t *testing.T) {
	encryptor, err := NewEncryptor("aes-128-cbc", "killer")
	if err != nil {
		fmt.Println(err)
		return
	}
	param, err := hex.DecodeString("ef2540d5f834cdff7724f34d22242ee0a7362af90a30f3f51f2f142af87b51cb")
	if err != nil {
		fmt.Println(err)
		return
	}
	result, err := encryptor.Decrypt(param)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(hex.EncodeToString(result))
}

func TestEncryptor_Encrypt_CBC(t *testing.T) {
	encryptor, err := NewEncryptor("aes-128-cbc", "killer")
	if err != nil {
		fmt.Println(err)
		return
	}

	result, err := encryptor.Encrypt(bytes.Repeat([]byte("a"),16))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(hex.EncodeToString(result))
}


func TestEncryptor_Decrypt_CFB(t *testing.T) {
	encryptor, err := NewEncryptor("aes-128-cfb", "killer")
	if err != nil {
		fmt.Println(err)
		return
	}
	param, err := hex.DecodeString("27fe5810f0d7109ecb66230d8760522342e705348223c74194783e5354b3a4ad")
	if err != nil {
		fmt.Println(err)
		return
	}
	result, err := encryptor.Decrypt(param)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(hex.EncodeToString(result))
}


func TestEncryptor_Decrypt_RC4(t *testing.T) {
	encryptor, err := NewEncryptor("rc4", "killer")
	if err != nil {
		fmt.Println(err)
		return
	}
	param, err := hex.DecodeString("a395e09871e16334fd9f9da33ffb67c0")
	if err != nil {
		fmt.Println(err)
		return
	}
	result, err := encryptor.Decrypt(param)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(hex.EncodeToString(result))

	encryptor, _ = NewEncryptor("rc4", "a2lsbGVybCPp7mlgedshQ+TWgabdtQ==")
	cleartext, _ := encryptor.Decrypt(bytesx.MustHexDecode("72193e8ab211df"))
	fmt.Printf("%s|%s\n",hex.EncodeToString(cleartext),string(cleartext))
	cleartext, _ = encryptor.Decrypt(bytesx.MustHexDecode("72193e8ab211df"))
	fmt.Printf("%s|%s\n",hex.EncodeToString(cleartext),string(cleartext))
	cleartext, _ = encryptor.Decrypt(bytesx.MustHexDecode("72193e8ab211df"))
	fmt.Printf("%s|%s\n",hex.EncodeToString(cleartext),string(cleartext))

}

func ExampleNewCBCDecrypter() {
	// Load your secret key from a safe place and reuse it across multiple
	// NewCipher calls. (Obviously don't use this example key for anything
	// real.) If you want to convert a passphrase to a key, use a suitable
	// package like bcrypt or scrypt.
	key, _ := hex.DecodeString("b36d331451a61eb2d76860e00c347396")
	ciphertext, _ := hex.DecodeString("ebaf33c9c39ac769642cf078360b0539d58f0ebbc168b03d185891d71879286d")

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// CBC mode always works in whole blocks.
	if len(ciphertext)%aes.BlockSize != 0 {
		panic("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(ciphertext, ciphertext)

	// If the original plaintext lengths are not a multiple of the block
	// size, padding would have to be added when encrypting, which would be
	// removed at this point. For an example, see
	// https://tools.ietf.org/html/rfc5246#section-6.2.3.2. However, it's
	// critical to note that ciphertexts must be authenticated (i.e. by
	// using crypto/hmac) before being decrypted in order to avoid creating
	// a padding oracle.

	fmt.Printf("%s\n", ciphertext)
	// Output: exampleplaintext
}

