package service

import (
	"net/http"
	"time"

	"github.com/rc452860/vnet/cmd/rpcx"

	"github.com/rc452860/vnet/cmd/rpcx/config"
	"github.com/rc452860/vnet/service"
	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
)

func HttpStart() {
	r := gin.Default()
	r.GET("/node_status", NodeStatus)
	r.GET("/node_users", NodeUserList)
	r.Run()
}

func NodeStatus(c *gin.Context) {
	token := c.Query("token")
	if token != viper.GetString(config.C_Token) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	ssservice := service.CurrentShadowsocksService()
	c.JSON(http.StatusOK, gin.H{
		"NodeId":         viper.GetString(config.C_NodeId),
		"Len":            len(ssservice.List()),
		"LastUpdateTime": time.Unix(g_updatetime, 0),
	})
}

func NodeUserList(c *gin.Context) {
	token := c.Query("token")
	if token != viper.GetString(config.C_Token) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	ssservice := service.CurrentShadowsocksService()
	var users []*rpcx.EnableUser = []*rpcx.EnableUser{}
	for _, item := range ssservice.List() {
		users = append(users, item.Data.(*rpcx.EnableUser))
	}
	c.JSON(http.StatusOK, users)
}
