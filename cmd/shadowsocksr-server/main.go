package main

import (
	"fmt"
	"github.com/rc452860/vnet/api/client"
	"github.com/rc452860/vnet/api/server"
	"github.com/rc452860/vnet/cmd/shadowsocksr-server/command"
	"github.com/rc452860/vnet/proxy"
	"github.com/rc452860/vnet/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)


func main(){
	logrus.SetLevel(logrus.DebugLevel)
	command.Execute(func() {
		client.SetHost(viper.GetString(command.API_HOST))
		nodeInfo := client.GetNodeInfo(viper.GetInt(command.NODE_ID),viper.GetString(command.KEY))
		logrus.WithFields(logrus.Fields{
			"nodeInfo":nodeInfo,
		}).Info("get node info success")
		service.ShadowsocksrServiceInstance.SetNodeInfo(&nodeInfo)
		users := client.GetUserList(viper.GetInt(command.NODE_ID),viper.GetString(command.KEY))
		logrus.WithFields(logrus.Fields{
			"firstLoadUserCount":len(users),
		}).Info("get user list success")
		// single port start
		if(nodeInfo.Single == 1){
			portStrArray :=strings.Split(nodeInfo.Port,",")
			ports := []int{}
			for _,item := range portStrArray{
				convertPort,err := strconv.Atoi(item)
				if err != nil{
					panic(fmt.Sprintf("port format error: %s",nodeInfo.Port))
				}
				ports = append(ports,convertPort)
			}

			for _,port := range ports{
				service.ShadowsocksrServiceInstance.ShadowsocksRProxy(viper.GetString(command.HOST),
					port,
					nodeInfo.Method,
					nodeInfo.Passwd,
					nodeInfo.Protocol,
					nodeInfo.ProtocolParam,
					nodeInfo.Obfs,
					nodeInfo.ObfsParam,
					&proxy.ShadowsocksRArgs{})

				 for _,item := range users{
					 service.ShadowsocksrServiceInstance.Shadowsocksrs[port].AddUser(item.Port,item.Passwd)
					 service.ShadowsocksrServiceInstance.AddUIDPortConvertItem(item.Uid,item.Port)
				 }
				err := service.ShadowsocksrServiceInstance.Shadowsocksrs[port].Start()
				if err !=nil{
					// TODO 错误处理
					panic(err)
				}
			}
		}else {
			// Multi Port start
			for _, item := range users {
				service.ShadowsocksrServiceInstance.ShadowsocksRProxy(viper.GetString(command.HOST),
					item.Port,
					nodeInfo.Method,
					item.Passwd,
					nodeInfo.Protocol,
					nodeInfo.ProtocolParam,
					nodeInfo.Obfs,
					nodeInfo.ObfsParam,
					&proxy.ShadowsocksRArgs{})

				service.ShadowsocksrServiceInstance.AddUIDPortConvertItem(item.Uid,item.Port)
				err := service.ShadowsocksrServiceInstance.Shadowsocksrs[item.Port].Start()
				if err != nil {
					// TODO error handle
					panic(err)
				}
			}
		}
		go timeWheelTask()
		go server.StartServer()
		// Wait for interrupt signal to gracefully shutdown the server with
		// a timeout of 5 seconds.
		quit := make(chan os.Signal)
		// kill (no param) default send syscall.SIGTERM
		// kill -2 is syscall.SIGINT
		// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
	})
}


func timeWheelTask(){
	timer := time.Tick(1*time.Second)
	tick := 0
	for  {
		<- timer
		if tick % 10 == 0{
			//updateUsers()

			traffic := service.ShadowsocksrServiceInstance.ReportTraffic()
			if len(traffic) > 0{
				client.PostAllUserTraffic(traffic,viper.GetInt(command.NODE_ID),viper.GetString(command.KEY))
			}
			online := service.ShadowsocksrServiceInstance.ReportOnline()
			if len(online) > 0{
				client.PostNodeOnline(online,viper.GetInt(command.NODE_ID),viper.GetString(command.KEY))
			}
			client.PostNodeStatus(service.ShadowsocksrServiceInstance.ReportNodeStatus(),viper.GetInt(command.NODE_ID),viper.GetString(command.KEY))
		}
		tick ++
	}
}

func updateUsers(){
	users := client.GetUserList(viper.GetInt(command.NODE_ID),viper.GetString(command.KEY))
	if service.ShadowsocksrServiceInstance.GetNodeInfo().Single == 1{
		for _,ssrserver := range service.ShadowsocksrServiceInstance.Shadowsocksrs{
			for _,user := range users{
				if user.Enable == 1{
					ssrserver.AddUser(user.Port,user.Passwd)
				}else{
					ssrserver.DelUser(user.Port)
				}
			}
		}
	}else {
		nodeInfo := service.ShadowsocksrServiceInstance.GetNodeInfo()
		for _,user := range users{
			if user.Enable == 1{
				service.ShadowsocksrServiceInstance.ShadowsocksRProxy(viper.GetString(command.HOST),
					user.Port,
					nodeInfo.Method,
					user.Passwd,
					nodeInfo.Protocol,
					nodeInfo.ProtocolParam,
					nodeInfo.Obfs,
					nodeInfo.ObfsParam,
					&proxy.ShadowsocksRArgs{})
				err := service.ShadowsocksrServiceInstance.Shadowsocksrs[user.Port].Start()
				if err != nil {
					// TODO error handle
					panic(err)
				}
			}else{
				service := service.ShadowsocksrServiceInstance.Shadowsocksrs[user.Port]
				if service == nil{
					continue
				}
				err := service.Close()
				if err != nil{
					panic(err)
				}
			}
		}
	}
}
