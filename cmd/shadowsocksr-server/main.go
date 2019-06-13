package main

import (
	"fmt"
	"github.com/rc452860/vnet/api"
	"github.com/rc452860/vnet/cmd/shadowsocksr-server/command"
	"github.com/rc452860/vnet/model"
	"github.com/rc452860/vnet/proxy"
	"github.com/rc452860/vnet/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strconv"
	"strings"
	"time"
)

var(
	nodeInfo model.NodeInfo
)

func main(){
	logrus.SetLevel(logrus.DebugLevel)
	command.Execute(func() {
		api.SetHost(viper.GetString(command.API_HOST))
		nodeInfo = api.GetNodeInfo(viper.GetInt(command.NODE_ID),viper.GetString(command.KEY))
		logrus.WithFields(logrus.Fields{
			"nodeInfo":nodeInfo,
		}).Info("get node info success")

		users := api.GetUserList(viper.GetInt(command.NODE_ID),viper.GetString(command.KEY))
		logrus.WithFields(logrus.Fields{
			"firstLoadUserCount":len(users),
		}).Info("get user list success")
		// 单端口模式启动
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
		timeWheelTask()
	})
}


func timeWheelTask(){
	timer := time.Tick(1*time.Second)
	tick := 0
	for  {
		<- timer
		if tick % 10 == 0{
			updateUsers()

			traffic := service.ShadowsocksrServiceInstance.ReportTraffic()
			if len(traffic) > 0{
				api.PostAllUserTraffic(traffic,viper.GetInt(command.NODE_ID),viper.GetString(command.KEY))
			}
			online := service.ShadowsocksrServiceInstance.ReportOnline()
			if len(online) > 0{
				api.PostNodeOnline(online,viper.GetInt(command.NODE_ID),viper.GetString(command.KEY))
			}
			api.PostNodeStatus(service.ShadowsocksrServiceInstance.ReportNodeStatus(),viper.GetInt(command.NODE_ID),viper.GetString(command.KEY))
		}
		tick ++
	}
}

func updateUsers(){
	users := api.GetUserList(viper.GetInt(command.NODE_ID),viper.GetString(command.KEY))
	if nodeInfo.Single == 1{
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
