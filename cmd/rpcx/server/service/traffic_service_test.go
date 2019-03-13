package service

import (
	"time"

	"github.com/rc452860/vnet/cmd/rpcx"
)

func ExampleRecordTrafficLog() {
	InitDBTest()
	tsi := GetTrafficServiceInstance()
	tsi.Start()
	TickDown(1*time.Second, 3, func() {
		tsi.RecordTrafficLogSync(&rpcx.PushUserTrafficRequest{
			TrafficInfo: []*rpcx.TrafficInfo{
				&rpcx.TrafficInfo{
					Uid:      1,
					Port:     1090,
					Upload:   100,
					Download: 200,
				},
				&rpcx.TrafficInfo{
					Uid:      4,
					Port:     10003,
					Upload:   200,
					Download: 400,
				},
			},
			NodeId: 2,
			Token:  "abc",
		})
	})
	tsi.Stop()
	time.Sleep(1 * time.Second)
	//Output:
}

func TickDown(d time.Duration, count int, callback func()) {
	tick := time.Tick(d)
	for i := 0; i < count; i++ {
		<-tick
		callback()
	}
}
