package api

import (
	"github.com/gin-gonic/gin"
)

func StartApi() {
	r := gin.Default()
	r.GET("/shadowsocks", ShadowsocksList)
	r.POST("/shadowsocks/add", ShadowsocksAdd)
	r.GET("/shadowsocks/get/:port", ShadowsocksGet)
	r.POST("/shadowsocks/start/:port", ShadowsocksStart)
	r.POST("/shadowsocks/stop/:port", ShadowsocksStop)
	r.DELETE("/shadowsocks/del/:port", ShadowsocksDel)
	r.GET("/system/info", SystemInfoGet)
	r.Run(":8080")
}
