package api

import (
	"github.com/gin-gonic/gin"
	"github.com/rc452860/vnet/service"
)

// SystemInfoGet return cpu mem network disk usage
func SystemInfoGet(g *gin.Context) {
	up, down := service.GetNetwork()
	g.JSON(200, NewSystemInfo(service.GetCPUUsage(),
		service.GetMemUsage(),
		service.GetDiskUsage(),
		up,
		down))
}
