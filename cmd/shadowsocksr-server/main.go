package main

import (
	"fmt"
	"github.com/rc452860/vnet/cmd/shadowsocksr-server/command"
	"time"
)

func main(){
	command.Execute(func() {
		fmt.Println("hello.")
		time.Sleep(3*time.Second)
		fmt.Println("world.")
	})
}

