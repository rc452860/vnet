package main

import (
	"fmt"
	"github.com/rc452860/vnet/network"
	"github.com/rc452860/vnet/proxy/server"
	"github.com/sirupsen/logrus"
)

func main() {
	//tls1.2_ticket_auth
	//auth_chain_a

	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&logrus.TextFormatter{})
	testSSR()

}

func testSS() {
	ss, _ := server.NewShadowsocks("0.0.0.0", "aes-128-cfb", "killer", 9090, server.ShadowsocksArgs{})
	ss.Start()
	input := ""
	_, _ = fmt.Scanf("%s", &input)
	fmt.Printf(input)
}

func testSSR() {
	ssr := &server.ShadowsocksRProxy{
		Host:          "0.0.0.0",
		Port:          9090,
		ProtocolParam: "",
		Protocol:      "auth_chain_a",
		Obfs:          "tls1.2_ticket_auth",
		ObfsParam:     "",
		Method:        "aes-128-cfb",
		Password:      "killer",
		Listener: &network.Listener{
			Addr: "0.0.0.0:9090",
		},
	}
	ssr.AddUser(1024, "killer")
	_ = ssr.StartTCP()
	input := ""
	_, _ = fmt.Scanf("%s", &input)
	fmt.Printf(input)
}
