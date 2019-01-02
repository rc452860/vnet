package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	listen, err := net.ListenPacket("udp", "0.0.0.0:8081")
	if err != nil {
		panic(err)
	}

	rage := time.Tick(100 * time.Millisecond)
	buf := make([]byte, 4096)
	for {
		<-rage
		_, _, err := listen.ReadFrom(buf)
		fmt.Printf("%v\n", buf[0:10])
		if err != nil {
			fmt.Print(err.Error())
			continue
		}

	}
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}
