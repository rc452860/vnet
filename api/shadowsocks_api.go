package api

import (
	"fmt"
	"strconv"

	"github.com/rc452860/vnet/service"

	"github.com/gin-gonic/gin"
)

type Shadowsocks struct {
	Host     string `form:"host" json:"host" binding:"required"`
	Method   string `form:"method" json:"method" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Port     int    `form:"port" json:"port" binding:"required"`
	Timeout  int64  `form:"timeout" json:"timeout"`
	Limit    string `form:"limit" json:"limit"`
}

func ShadowsocksAdd(c *gin.Context) {
	var (
		ss       Shadowsocks
		response Response
	)
	if err := c.ShouldBindJSON(&ss); err != nil {
		response.Code = ERROR
		response.Message = err.Error()
		c.JSON(500, response)
		return
	}
	err := service.CurrentShadowsocksService().Add(
		ss.Host,
		ss.Method,
		ss.Password,
		ss.Port,
		ss.Limit,
		ss.Timeout,
	)
	if err != nil {
		c.JSON(500, Error(err))
		return
	}
	c.JSON(200, Success(nil))
}

func ShadowsocksGet(c *gin.Context) {
	port, err := strconv.Atoi(c.Param("port"))
	if err != nil {
		c.JSON(500, Error(err))
		return
	}
	service := service.CurrentShadowsocksService().Get(port)
	if service != nil {
		c.JSON(200, service)
		return
	} else {
		c.JSON(200, Failure("port is not found"))
		return
	}
}

func ShadowsocksList(c *gin.Context) {
	c.JSON(200, service.CurrentShadowsocksService().List())
}

func ShadowsocksStart(c *gin.Context) {
	port, err := strconv.Atoi(c.Param("port"))
	if err != nil {
		c.JSON(500, Error(err))
		return
	}
	err = service.CurrentShadowsocksService().Start(port)
	if err != nil {
		c.JSON(500, Error(err))
		return
	}
	c.JSON(200, Success(fmt.Sprintf("start %v success", port)))
}

func ShadowsocksStop(c *gin.Context) {
	port, err := strconv.Atoi(c.Param("port"))
	if err != nil {
		c.JSON(500, Error(err))
		return
	}
	err = service.CurrentShadowsocksService().Stop(port)
	if err != nil {
		c.JSON(500, Error(err))
		return
	}
	c.JSON(200, Success(fmt.Sprintf("stop %v success", port)))
}

func ShadowsocksDel(c *gin.Context) {
	port, err := strconv.Atoi(c.Param("port"))
	if err != nil {
		c.JSON(500, Error(err))
		return
	}
	if err := service.CurrentShadowsocksService().Del(port); err != nil {
		c.JSON(500, Error(err))
		return
	}
	c.JSON(200, Success(fmt.Sprintf("del %v success", port)))
}
