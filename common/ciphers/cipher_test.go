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
	ciphertext := bytesx.MustHexDecode("7ae0cd7ede4f99f5f51d0d52cd3026c3bb266241305dc4f8f968fe55703006fcd456724001c7dcebd4f5af51df89ab9921173a0f60735922d3ce52f58e8fb0aceb65fa834425eb3826961e5034bca64d1d7ccb8485cea8cfcdb011f3cef766f8275f83609cd844b658877ed79cd3612f26f6449494d660dc9039175e627af2d224363e8262796738064fdf3e53e42546d5b3cf32bea66dac7015b6c87967d2fd954d85d8cf8dfaa70330dbd98029e86ae30ac820f60496ab0a5e8a357c634183d640e9362445eaaec9e3d682c037889b539e1cedf9843cd152bd169c0d2fc55e88747773434240cb0686735251eb46317e7e31b504a253d9e2d15ef8ec0e61f515ec7a3d4ed57691e1fc0f109d9415919a7e2d642b188c522facae4bbf8853e1a02fe060d1a14b33b2e5a24b8ef74f30053feb45603fee1eece2416d2d457dd3f24cb1d5143a7a523f20e4c75912ae1c8774595f64fac73031010677137f9a9306b4f40dbdcdf997a52fdea60c863f71cc88ffaafa8bb2a561356e05fff673f72f48088303e8f39a03b0e1b7ab0423a4974209f996b6c0fbd9c885a3e36899873828e0a9b9bfda795f57fa4511edd6371f42e8fddb7c11df3ae2c167f27d53550532a5847b4d90fbb680f4a5931f503ebf866bbc760f2fb377aa1a9ac81df09c60a84a70c4699c2c9c7ecbce3ec834c9b9ef7ccdd9adc41f693f4c27d048c9ccbef447fef98e6350a9934174b327ed92167b3a15718581e5e0ea5701e13135c245ca50be11753f54538106d649f70a3044940d80d455e0b00696931454f0e5e5f135e5cf37d332f07a8907fc592da8adf5d9056a506dd819e556802c19487ae01a2c276754af3012b3dccf8e53fce4810faf4d56c5aeb0f2d114b18a025527b5504a655c3268c60ce21f8c")
	cleartext, err := encryptor.Decrypt(ciphertext)
	fmt.Println(hex.EncodeToString(cleartext))
	a, _ := aes.NewCipher(encryptor.Key)
	cfb := cipher.NewCFBDecrypter(a, ciphertext[:16])
	buf := make([]byte, len(ciphertext))
	cfb.XORKeyStream(buf, ciphertext[16:])
	fmt.Println(hex.EncodeToString(buf))

	//raw := []byte("abc")
	//fmt.Println(hex.EncodeToString(raw))
	//result, err := encryptor.Encrypt(raw)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println("ciphertext result: " + hex.EncodeToString(result))
	//result, err = encryptor.Decrypt(result)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println("plain text: " + hex.EncodeToString(result))
	//if hex.EncodeToString(result) != "616263" {
	//	t.Fail()
	//}
	//Output:
	// 142666ee03da64202914e3a3548b948470e50542b9d96a58e83647af175f8313902d2098b5230cb8f6d30caa1eca87c0e3f255b9baf72d459b061292a41f92285047dbd012c3585b2ddae3508491f29246cae1dbc61d591122cf8b5fdb7a14d741e35ad7efb6ad1f99fa8ce014dcc00914804bda3eb66233edc9ad4338f87fb2c50b0fa1371aa38f1dc141cb9079fc450e814b5b2f3a998a8026ea250969e7a07ada1a08cd1fd1f9dcbcbc7c7597a274cd6735b6ac23838b72852010307d86b3daa150b65ef478eeaf527e737d1c15e9837558a13330f456578f9647ffc0c95f3508a2d346ec16d9ef46065aa5100d5570368dba24e580a9adfa5fcf5d31efef881f3d747a41a3965f7fb1ed2999b1468140d6c9330e3b5088f8264c8d467bbd5717d286d298db2f92f65ba433927f1348184d1d5a26bb79e3eceb811a57091804f39512c72b73e864a6363db3961df579130375f95d62315c58a2e5ad59153bed439f3ffd6d4c1ca464b8b4ef7ec64ccff6080653931ea79164b56e898d3b0a91fb5070f7ddf80bc6e7d194a4ae41de925516b6510ee75924f8ae99b7de931629ad23f42d5655f2fdfd406cc28a02e222d7db337784b98213f74bfe2f5f53e853ba6bf68b4db03c0023b954ce7a969af5311ee7169783b0686cc773e8ebf6aed506955e2ca25808ab803c9e1dfbcc76df02609ea4f4a0e0360a56c39da95819071e0b56d4279f3254cb65b39860b10863e974e5f21cc702ea6c7932a9601bab76618e7bea75a1b8632e4ad98eec47d6584d60b4347226eb57096c67ae6c88f3afba3d37d1361eed47bc81c4d1ee73fd25a3d95b006776033aa323e0093909f7a55729d4af1004222f726d78f1621969956e420b922cf34843e747
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

	result, err := encryptor.Encrypt(bytes.Repeat([]byte("a"), 16))
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
	fmt.Printf("%s|%s\n", hex.EncodeToString(cleartext), string(cleartext))
	cleartext, _ = encryptor.Decrypt(bytesx.MustHexDecode("72193e8ab211df"))
	fmt.Printf("%s|%s\n", hex.EncodeToString(cleartext), string(cleartext))
	cleartext, _ = encryptor.Decrypt(bytesx.MustHexDecode("72193e8ab211df"))
	fmt.Printf("%s|%s\n", hex.EncodeToString(cleartext), string(cleartext))

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
