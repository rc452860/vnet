package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rc452860/vnet/config"

	"github.com/rc452860/vnet/log"
	"github.com/rc452860/vnet/proxy/server"
)

var logging *log.Logging

func init() {
	logging = log.GetLogger("root")
	logging.Level = log.INFO
}

func main() {
	config.LoadConfig("config.json")
	shadowsocks := server.NewShadowsocsk("0.0.0.0", "aes-128-gcm", "killer", 1090)
	if err := shadowsocks.Start(); err != nil {
		logging.Err(err)
		return
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}
