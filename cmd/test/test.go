package main

import (
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"runtime/pprof"

	"github.com/rc452860/vnet/api"
	"github.com/rc452860/vnet/log"
)

func main() {
	f, err := os.Create("profile.pprof")
	if err != nil {
		log.Err(err)
		return
	}
	pprof.StartCPUProfile(f)
	go api.StartApi()
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	mem, err := os.Create("mem.pprof")
	if err != nil {
		log.Err(err)
	}
	runtime.GC()
	pprof.WriteHeapProfile(mem)
	mem.Close()
	pprof.StopCPUProfile()
	f.Close()
}
