package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	for {
		data := <-c
		fmt.Println(data.String())
		if data == os.Interrupt {
			return
		}
	}

}
