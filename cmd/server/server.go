package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rc452860/vnet/config"
	"github.com/rc452860/vnet/log"
	"github.com/rc452860/vnet/proxy/server"
	"github.com/rc452860/vnet/utils/datasize"
)

var logging *log.Logging

func init() {
	logging = log.GetLogger("root")
	logging.Level = log.INFO
}

func main() {
	_, err := config.LoadConfig("config.json")
	if err != nil {
		logging.Err(err)
		return
	}

	host := flag.String("host", "0.0.0.0", "shadowsocks server host")
	method := flag.String("method", "aes-128-gcm", "shadowsocks method")
	password := flag.String("password", "killer", "shadowsocks password")
	port := flag.Int("port", 1090, "shadowsocks port")
	limit := flag.String("limit", "", "shadowsocks traffic limit exp:4MB")
	flag.Parse()

	shadowsocks, err := server.NewShadowsocks(*host, *method, *password, *port, *limit, 0)
	if err != nil {
		logging.Err(err)
		return
	}
	if err := shadowsocks.Start(); err != nil {
		logging.Err(err)
		return
	}
	go func() {
		tick := time.Tick(1 * time.Second)
		for {
			<-tick
			upSpeed, _ := datasize.HumanSize(shadowsocks.UpSpeed)
			downSpeed, _ := datasize.HumanSize(shadowsocks.DownSpeed)
			upBytes, _ := datasize.HumanSize(shadowsocks.UpBytes)
			downBytes, _ := datasize.HumanSize(shadowsocks.DownBytes)
			logging.Info("[upspeed: %s] [downspeed: %s] [up: %s] [down: %s]", upSpeed, downSpeed, upBytes, downBytes)
		}
	}()
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}
