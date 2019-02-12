// Package socks implements essential parts of SOCKS protocol.
package socks

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSocks5Addr_GetRaw(t *testing.T) {
	tests := []struct {
		name    string
		ss      *Socks5Addr
		wantRaw []byte
		wantErr bool
	}{
		{
			"aaa",
			NewSSProtocol(AtypIPv4, 3306, "127.0.0.1"),
			NewSSProtocol(AtypIPv4, 3306, "127.0.0.1").MustGetRaw(),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRaw, err := tt.ss.GetRaw()
			if (err != nil) != tt.wantErr {
				t.Errorf("Socks5Addr.GetRaw() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRaw, tt.wantRaw) {
				t.Errorf("Socks5Addr.GetRaw() = %v, want %v", gotRaw, tt.wantRaw)
			}
		})
	}
}

func ExampleSocks5Addr_GetRaw() {
	fmt.Printf("%v\n", NewSSProtocol(AtypIPv4, 3306, "127.0.0.1").MustGetRaw())
	ss := SplitAddr(NewSSProtocol(AtypIPv4, 3306, "127.0.0.1").MustGetRaw())
	if ss == nil {
		fmt.Println("ss is null")
	}
	fmt.Printf("%v\n", ss.MustGetRaw())

	fmt.Printf("%v\n", NewSSProtocol(AtypDomainName, 3306, "baidu.com").MustGetRaw())
	ss = SplitAddr(NewSSProtocol(AtypDomainName, 3306, "baidu.com").MustGetRaw())
	if ss == nil {
		fmt.Println("ss is null")
	}
	fmt.Printf("%v\n", ss.MustGetRaw())
	//Output:
	//[1 127 0 0 1 12 234]
	//[1 127 0 0 1 12 234]
	//[3 9 98 97 105 100 117 46 99 111 109 12 234]
	//[3 9 98 97 105 100 117 46 99 111 109 12 234]
}
