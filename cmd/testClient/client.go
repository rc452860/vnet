package main

import (
	_ "net/http/pprof"
	"os"
	"os/signal"

	thttp "github.com/rc452860/vnet/testing/servers/http"
)

func main() {
	thttp.StartFakeFileServer()
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Kill, os.Interrupt)
	<-ch
}
