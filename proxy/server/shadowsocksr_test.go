package server

import "fmt"

func ExampleShadowsocksR(){
	ssr := &ShadowsocksRProxy{
		Host:"0.0.0.0",
		Port:8838,
		ProtocolParam:"",
		Protocol:"auth_chan_a",
		Obfs:"tls1.2_ticket_auth",
		ObfsParam:"",
	}
	_ = ssr.StartTCP()
	_,_ = fmt.Scanln(nil)
	//Output:
}
