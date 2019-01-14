package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rc452860/vnet/ciphers"
	"github.com/rc452860/vnet/proxy/server"
)

func main() {
	fmt.Printf("%v", ciphers.GetSupportCiphers())
	ss, _ := server.NewShadowsocks("0.0.0.0", "chacha20", "killer", 1090, "4MB", 0)
	ss.Start()

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

}
