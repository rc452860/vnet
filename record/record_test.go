package record

import (
	"fmt"
	"testing"
	"time"

	"github.com/rc452860/vnet/common/eventbus"
	"github.com/rc452860/vnet/utils/addr"
)

func Test_GetLastOneMinuteOnlineByPort(t *testing.T) {
	GetGRMInstanceWithTick(100 * time.Millisecond)
	for i := 0; i < 10000000; i++ {

		eventbus.GetEventBus().Publish("record:proxyRequest", ConnectionProxyRequest{
			ConnectionPair: ConnectionPair{
				ProxyAddr:  addr.ParseAddrFromString("tcp", fmt.Sprintf("0.0.0.0:%v", i%100+100)),
				ClientAddr: addr.ParseAddrFromString("tcp", fmt.Sprintf("192.168.1.%v:%v", i%255, 35)),
				TargetAddr: addr.ParseAddrFromString("tcp", fmt.Sprintf("192.168.1.%v:%v", i%255, 35)),
			},
		})
		// time.Sleep(5 * time.Millisecond)
	}

	// if addr.GetIPFromAddr(instance.GetLastOneMinuteOnlineByPort()[8080][0]) != "192.168.1.1" {
	// 	t.FailNow()
	// }
	t.Log("aaa")
}
