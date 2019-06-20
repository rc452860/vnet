package main

import (
	"github.com/rc452860/vnet/api/client"
	"github.com/rc452860/vnet/api/server"
	"github.com/rc452860/vnet/cmd/shadowsocksr-server/command"
	"github.com/rc452860/vnet/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
	"time"
)


func main(){
	logrus.SetLevel(logrus.InfoLevel)
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
		for i := 0;i < len(users);i++{
			if err := service.ShadowsocksrServiceInstance.AddUser(&users[i]);err != nil{
				logrus.Error(err)
				return
			}
		}

		go timeWheelTask()
		go server.StartServer(nodeInfo.PushPort)
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
			traffic := service.ShadowsocksrServiceInstance.ReportTraffic()
			if len(traffic) > 0{
				logrus.Infof("prepare report traffic data, data length: %v",len(traffic))
				client.PostAllUserTraffic(traffic,viper.GetInt(command.NODE_ID),viper.GetString(command.KEY))
			}
			online := service.ShadowsocksrServiceInstance.ReportOnline()
			if len(online) > 0{
				logrus.Infof("prepare report online data, data length: %v",len(online))
				client.PostNodeOnline(online,viper.GetInt(command.NODE_ID),viper.GetString(command.KEY))
			}
			client.PostNodeStatus(service.ShadowsocksrServiceInstance.ReportNodeStatus(),viper.GetInt(command.NODE_ID),viper.GetString(command.KEY))
		}
		tick ++
	}
}

