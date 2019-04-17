package obfs

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
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

	//Output:
	//&{72340172838076673 72340172838076673}
	//&{11655686789823302041 3472380973156407715}
}


func ExampleHashMd5(){
	md5Data := md5.Sum([]byte("aaa"))
	fmt.Println(hex.EncodeToString(md5Data[:]))

	//Output:
	//47bce5c74f589f4867dbd57e9ca9f808
}