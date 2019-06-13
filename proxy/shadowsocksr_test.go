package proxy

import (
	"fmt"
	"github.com/rc452860/vnet/common/network/ciphers"
	"github.com/rc452860/vnet/utils/socksproxy"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

func ExampleShadowsocksR() {
	host := "127.0.0.1"
	port := 9090
	method := "aes-128-cfb"
	password := "killer"

	transport := &http.Transport{
		Proxy: nil,
		Dial: func(network, addr string) (net.Conn, error) {
			con, _ := net.Dial("tcp", fmt.Sprintf("%s:%v", host, port))

			c, err := ciphers.CipherDecorate(password, method, con)
			_,_ = c.Write(socksproxy.ParseAddr(addr).Raw)
			return c, err
		},
		TLSHandshakeTimeout: 10 * time.Second,
	}
	client := &http.Client{
		Transport: transport,
	}
	response, err := client.Get("httpserver://163.com")
	if err != nil {
		return
	}
	if response.StatusCode < 200 && response.StatusCode > 400 {
		fmt.Println("httpserver status error")
	}

	text, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(text))
	//Output:
}
