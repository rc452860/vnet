package main

import (
	"crypto/rand"
	"io"
	"net"
)

func main() {
	listen, err := net.ListenPacket("udp", "0.0.0.0:8082")
	if err != nil {
		panic(err)
	}
	dstAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:8081")
	if err != nil {
		panic(err)
	}
	tmp := make([]byte, 4096)
	if _, err = io.ReadFull(rand.Reader, tmp); err != nil {
		panic(err)
	}

	for i := 0; i < 100000; i++ {
		_, err := listen.WriteTo(tmp, dstAddr)
		if err != nil {
			panic(err)
		}
	}

}
