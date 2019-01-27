package service

import (
	"runtime"
	"testing"
	"time"

	"github.com/rc452860/vnet/common/log"
)

func TestMemProblem(t *testing.T) {
	log.GetLogger("root").Level = log.INFO
	for i := 10000; i < 20000; i++ {
		CurrentShadowsocksService().Add("0.0.0.0", "aes-128-cfb", "killer", i, "", 0)
		CurrentShadowsocksService().Start(i)
		CurrentShadowsocksService().Stop(i)
		CurrentShadowsocksService().Del(i)
	}
	time.Sleep(3 * time.Second)
	runtime.GC()

}
