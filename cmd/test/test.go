package main

import (
	"fmt"
	"net"
)

type Test struct {
	Name string
}

func main() {
	listen, _ := net.Listen("tcp", "0.0.0.0:8080")
	fmt.Print(listen.Addr().String())
}
