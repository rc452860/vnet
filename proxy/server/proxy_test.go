package server

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/rc452860/vnet/record"
)

func Test_LastOneMinute(t *testing.T) {
	proxy := NewProxyServiceWithTick(50 * time.Millisecond)
	proxy.Start()
	proxyAddr, _ := net.ResolveTCPAddr("tcp", "0.0.0.0:8000")
	targetAddr, _ := net.ResolveTCPAddr("tcp", "0.0.0.0:8081")

	for i := 0; i < 10; i++ {
		client, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf("1.1.1.1:%v", i+1000))
		proxy.MessageRoute <- record.ConnectionProxyRequest{
			ClientAddr:   client,
			ProxyAddr:    proxyAddr,
			TargetAddr:   targetAddr,
			TargetDomain: "baidu.com",
		}
		time.Sleep(10 * time.Millisecond)
	}

	var r int
	proxy.LastOneMinuteConnections.Range(func(key, value interface{}) {
		if result, ok := value.([]record.ConnectionProxyRequest); ok {
			// t.Log(key)
			for _, item := range result {
				t.Log(item.ClientAddr.String())
			}
			r++
		}
	})
	if r > 10 {
		t.FailNow()
	}
}
