package main

import (
	"flag"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/rc452860/vnet/service"

	"github.com/AlecAivazis/survey"
	"github.com/rc452860/vnet/config"
	"github.com/rc452860/vnet/db"
	"github.com/rc452860/vnet/log"
	"github.com/rc452860/vnet/proxy/server"
	"github.com/rc452860/vnet/utils/datasize"
)

func main() {
	conf, err := config.LoadDefault()
	log.Info("cpu core: %d", runtime.NumCPU())
	runtime.GOMAXPROCS(runtime.NumCPU())
	if err != nil {
		log.Err(err)
		return
	}

	if conf.Mode == "db" {
		DbStarted()
	} else {
		BareStarted()
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}

func DbStarted() {
	conf := config.CurrentConfig()
	if conf.DbConfig.Host == "" {
		survey.AskOne(&survey.Input{
			Message: "what is your host address?",
		}, &conf.DbConfig.Host, nil)
	}

	if conf.DbConfig.Port == "" {
		survey.AskOne(&survey.Input{
			Message: "what is your database port?",
		}, &conf.DbConfig.Port, nil)
	}

	if conf.DbConfig.User == "" {
		survey.AskOne(&survey.Input{
			Message: "what is your username?",
		}, &conf.DbConfig.User, nil)
	}

	if conf.DbConfig.Passwd == "" {
		survey.AskOne(&survey.Password{
			Message: "what is your password?",
		}, &conf.DbConfig.Passwd, nil)
	}

	if conf.DbConfig.Database == "" {
		survey.AskOne(&survey.Input{
			Message: "what is your database name?",
		}, &conf.DbConfig.Database, nil)
	}

	tick := time.Tick(time.Second)
	for {
		config.SaveConfig()
		users, err := db.GetEnableUser()
		if err != nil {
			log.Err(err)
			return
		}

		sslist := service.CurrentShadowsocksService().List()
		for _, user := range users {
			flag := false
			for _, ss := range sslist {
				if ss.Port == user.Port {
					flag = true
				}
			}
			if !flag {
				log.Info("start port [%d]", user.Port)
				StartShadowsocks(user)
			}
		}

		for _, ss := range sslist {
			flag := false
			for _, user := range users {
				if ss.Port == user.Port {
					flag = true
				}
			}
			if !flag {
				log.Info("stop port [%d]", ss.Port)
				service.CurrentShadowsocksService().Del(ss.Port)
			}
		}
		//TODO update user upload and download
		<-tick

	}
}

func StartShadowsocks(user db.User) {
	limit, err := datasize.HumanSize(user.Limit)
	if err != nil {
		log.Error("limit: %d, error info:%s", user.Limit, err.Error())
		return
	}
	err = service.CurrentShadowsocksService().Add("0.0.0.0",
		user.Method,
		user.Password,
		user.Port,
		limit,
		3*time.Second)
	if err != nil {
		log.Info("[%d] add failure, case %s", user.Port, err.Error())
	}
	err = service.CurrentShadowsocksService().Start(user.Port)
	if err != nil {
		log.Info("[%d] started failure, case: %s", user.Port, err.Error())
	}
}

func BareStarted() {
	host := flag.String("host", "0.0.0.0", "shadowsocks server host")
	method := flag.String("method", "aes-128-cfb", "shadowsocks method")
	password := flag.String("password", "killer", "shadowsocks password")
	port := flag.Int("port", 1090, "shadowsocks port")
	limit := flag.String("limit", "", "shadowsocks traffic limit exp:4MB")
	flag.Parse()
	shadowsocks, err := server.NewShadowsocks(*host, *method, *password, *port, *limit, 0)
	if err != nil {
		log.Err(err)
		return
	}
	if err := shadowsocks.Start(); err != nil {
		log.Err(err)
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
			log.Info("[upspeed: %s] [downspeed: %s] [up: %s] [down: %s]", upSpeed, downSpeed, upBytes, downBytes)
		}
	}()

}
