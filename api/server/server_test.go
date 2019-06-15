package server

import (
	"github.com/gin-gonic/gin"
	"os"
	"os/signal"
	"syscall"
)

func ExampleStartServer() {
	gin.SetMode(gin.ReleaseMode)
	StartServer()
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	<-c
	//Output:
}
