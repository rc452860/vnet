package api

import (
	"github.com/labstack/echo"
	"github.com/rc452860/vnet/log"
)

var Router *echo.Echo
var logging *log.Logging

func init() {
	logging := logging
	Router = echo.New()

}

func InitRouter() {
	Router.POST("/shadowsocks/add", ShadowsocksAdd)
	Router.GET("/shadowsocks/list", ShadowsocksList)
	Router.GET("/shadowsocks/status", ShadowsocksStatus)
	Router.Post("/shadowsocks/restart", ShadowsocksRestart)
}
