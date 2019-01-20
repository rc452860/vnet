package record

import (
	"fmt"
	"testing"
	"time"

	"github.com/rc452860/vnet/comm/eventbus"
	"github.com/rc452860/vnet/utils/addr"
)

func Test_GetLastOneMinuteOnlineByPort(t *testing.T) {
	instance := GetGRMInstanceWithTick(100 * time.Millisecond)
	for i := 0; i < 10; i++ {

		eventbus.GetEventBus().Publish("record:proxyRequest", ConnectionProxyRequest{
			ConnectionPair: ConnectionPair{
				ProxyAddr:  addr.ParseAddrFromString("tcp", "0.0.0.0:8080"),
				ClientAddr: addr.ParseAddrFromString("tcp", fmt.Sprintf("192.168.1.1:%v", i)),
				TargetAddr: addr.ParseAddrFromString("tcp", fmt.Sprintf("192.168.1.2:%v", i)),
			},
		})
		time.Sleep(5 * time.Millisecond)
	}

	if instance.GetLastOneMinuteOnlineByPort()[8080][0] != "192.168.1.1" {
		t.FailNow()
	}
}
