package obfs

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/rc452860/vnet/utils/binaryx"
	"net"
	"testing"
)

func ExampleAuthChain() {
	var a Plain = NewAuthBase("acda")
	fmt.Println(a.GetMethod())
	s := NewServerInfo()
	a.SetServerInfo(s)
	s.SetHost("baidu.com")
	fmt.Println(a.GetServerInfo().GetHost())
	//Output:
}

func ExampleXorShift128Plus() {
	a := NewXorShift128Plus()
	a.InitFromBin(bytes.Repeat([]byte{byte(0x01)}, 16))
	fmt.Println(a)
	a = NewXorShift128Plus()
	a.InitFromBinLen(bytes.Repeat([]byte{byte(0x01)}, 16), 2)
	fmt.Println(a)
	fmt.Printf("%v\n", a.Next())

	a = NewXorShift128Plus()
	a.InitFromBinLen(MustHexDecode("1c61777508c444705f7ec9092de53b7e"), 5)
	fmt.Printf("%v\n", a.Next())

	//Output:
	//&{72340172838076673 72340172838076673}
	//&{11655686789823302041 3472380973156407715}
	//5864189510468454840

}

func ExampleHashMd5() {
	md5Data := md5.Sum([]byte("aaa"))
	fmt.Println(hex.EncodeToString(md5Data[:]))

	//Output:
	//47bce5c74f589f4867dbd57e9ca9f808
}

func ExampleAuthChainA() {
	server := GetAuth()
	testData, _ := hex.DecodeString("1f1f645b1c61777508c444705f7ac909ce549d8859f4959996b9b817cfccf89560c0716f3e7ec8ad809fea1785b0f2ae162789d30cec6a589cb71d75a796584d4e085b15189bee7f1daff841d802a70e36dc9a541937c38cbf2cb8252fc06f5f2a374bd12454316d7265921e84d8e662749a1fa93198dc3cc946a9aafe0b30b2c013a5eaa2bf9b9e74aee69d7996b62e9e3302e87cefe9137c0d4572e50944e1c6a7856549bf1ef4b10ee1b2c370420c53e9e7e78f975d022257236da092e424c96746d88d28b764ddc45d4f1b061b94424cc02b4d6b4a8ce0f97d30235a4dfa0fb1acd51ade9fe78be6991c5c2d6caad8028cf7a0660ca3cf0f805c42b4c59ad1628813df92866b9007000fc017b4995e087aa35e0d77a6b2e7d934672b37a925cca50ed9ff715aab68a9b812afa93228bf54e0ded51d989bdd9aa86a52a05a9283e955dd34f81af06d43650a43651ebdd5e36c9b42697f636d1275b2d61389a62c3cafed92e532316ce4d3d2e913c81d6f225192da176cfa007307f11b605918240c5d3ebed516f3c1ea657e6c6337f1bd1cb34e01d34660f592408cf50b936f9beec5409f8718ac4e3f833bb400cb01ceec7d67fa34b7971167f56dc8d335d9865d11e52d78c773a99ce2755fa10d537556252cc0604640bb63d98140edfa2def15d316e183dab62e393aa7a3b885a2eb25e0d2476b6a008c1987772898ad4777b75a36b22a8c52ebb248cf9db6916ad622e3bd049c0b1acf8610d5f1ce94e8b47df03e2168b11d067916ce6465501f4844fc5c8cec9e95455703c1330ce62045fd3ff0ef6d3107a7c7d5c01ff3f48dae7582a1c80779c36a433f24d919d4b1528e7cac85790df311c4231907e5989f11053603d53952b7d8ef815ac831a3a4e37e467327b50042f087a3bce89f347b98e9cdbf7f395b00d43d384455ae8fbcc98ddebbddaac8a0f4e784916ee12c4abd2a7fc3f792d39b242b57c67ebe22867c92460bbe15450eca27aa1c8a0df86522e050bda1f8a3d80352f5e4017d8b1e06e0df6d0f230b4aefdb7d7ab517725d347582c766963a7f30759f9194ba431decf37f422c2716ac4dfcf3cb1785243012a095b3c3276329a3f12effd50a2696eb79ceec8b245a7605b84b76e0f9ac1eeb1694e723285cc16ba613f8407905090aa0088b1b18e16abb6cc6ce22f8cb8a6858546cbe18eb277df693be995e3b2a61a4278820cba983e64b9397444edc239ee5bd06bb69d370383cb13f5d23cca1b078996a91550e74578b2215c2aeaf5a6a00267f878d58860fbc415ce3b9239f93ad97acfb9d49e3eca3a5b9832b773224ab54ba60af74e098465fbee992f4639cb6dcd43b648441ae89dc8039b9d72a8af09ca33f2d5b3056e4d193ae8488092b95573ac84e49542b63c75ca9e982109b56b7f1850974f5b7ad0eb3fe9e856e6513801608fdc9cc125fc35cc23faf17c265ff05f18676659bef8082fcc6b78fd07d83b765f0b11333af84498d27675fd3e6de77d8b1fe3c630fe36f538535adc8f6cbf8ae3e6671b81b59fb14830d5d72e008d41e032f8aefc3f3a588a7e47f86ad479912d11a4b96d840a316837b082a2bd50f8322ebee431c1bd20f2ddbe9b8f59d749c8da29f90d4f28314e9fea064292e049cbb6d7e7d3aa389c6343af69710f4a3be10913c21a10c2922564f6d1a66be2366a95697515b4bf2291e427aa46eb81c5aa4e71022d3aaab458d556b13dd632ec6a2bf311b4a085fd91210415c7c0ce3b0e03de498bd41a0b046efcf91b8654370b302b509acc995a3c99a700fdb50d9eed0260d62a8568c837f10d06865cd7ddcbf57ca3f2c3beac7137c5a6d22aba73fcfbd03b16b6a4407026b46e7c1cad47c1064cfac076c")
	result, sendback, err := server.ServerPostDecrypt(testData)
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Printf("sendback: %v,result: %s\n", sendback, result)
	//Output:
	//
}

func ExampleAuthChainAClient() {
	client := GetAuth()
	server := GetAuth()
	ciphertext, err := client.ClientPreEncrypt([]byte("hello"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(hex.EncodeToString(ciphertext))
	result, sendback, err := server.ServerPostDecrypt(ciphertext)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("sendback: %v,result: %s\n", sendback, result)

	result, err = server.ServerPreEncrypt([]byte("hello"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(hex.EncodeToString(result))

	result, err = client.ClientPostDecrypt(result)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s|%s", string(result), hex.EncodeToString(result))
	//Output:
}

func TestTestAuthChainUDP(t *testing.T) {
	client := GetAuth()
	server := GetAuth()
	result, err := client.ClientUDPPreEncrypt([]byte("hello"))
	if err != nil {
		fmt.Println(err)
		return
	}
	result, uid, err := server.ServerUDPPostDecrypt(result)
	if err != nil {
		fmt.Println(err)
		return
	}
	if string(result) != "hello" && binaryx.LEBytesToUInt32([]byte(uid)) != 1024 {
		t.Fatal("service decrypt result is not equal hello")
		return
	}
	result, err = server.ServerUDPPreEncrypt([]byte("hello"), []byte(uid))
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Printf("service encrypt# %s | %s\n",string(result),hex.EncodeToString(result))

	result, err = client.ClientUDPPostDecrypt(result)
	if err != nil {
		fmt.Println(err)
		return
	}
	if string(result) != "hello" {
		t.Fatal("service decrypt result is not equal hello")
		return
	}

	//Output:
}
func UpdateUser(uid []byte) {
	fmt.Printf("uid:%v \n", binaryx.LEBytesToUInt32(uid))
}

func GetAuth() Plain {
	authChainA := NewAuthChainA("auth_chain_a")
	serverInfo := NewServerInfo()
	serverInfo.GetUsers()[string(binaryx.LEUint32ToBytes(1024))] = "killer"
	serverInfo.SetClient(net.ParseIP("127.0.0.1"))
	serverInfo.SetPort(8080)
	serverInfo.SetProtocolParam("1024:killer")
	serverInfo.SetIv(MustHexDecode("271d7f17d03ed7cd1f44327456aebfa2"))
	serverInfo.SetRecvIv(MustHexDecode("271d7f17d03ed7cd1f44327456aebfa2"))
	serverInfo.SetKeyStr("killer")
	serverInfo.SetKey(MustHexDecode("b36d331451a61eb2d76860e00c347396"))
	serverInfo.SetHeadLen(30)
	serverInfo.SetTCPMss(1460)
	serverInfo.SetBufferSize(32*1024 - 5 - 4)
	serverInfo.SetOverhead(9)
	serverInfo.SetUpdateUserFunc(UpdateUser)
	authChainA.SetServerInfo(serverInfo)
	return authChainA
}
