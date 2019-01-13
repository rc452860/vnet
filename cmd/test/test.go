package main

import (
	"fmt"
	"time"

	"github.com/rc452860/vnet/service"
)

func main() {

	tick := time.Tick(time.Second)
	for {
		<-tick
		up, down := service.GetNetwork()
		fmt.Printf("cpu: %v, mem: %v, network: %v - %v, disk: %v\n",
			service.GetCPUUsage(),
			service.GetMemUsage(),
			up, down,
			service.GetDiskUsage())
	}

	// stop := make(chan os.Signal)
	// signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// <-stop
	// if err != nil {
	// 	log.Err(err)
	// }

}
