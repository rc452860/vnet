package record

import (
	"context"
	"net"
	"sync"

	"github.com/rc452860/vnet/comm/eventbus"

	"github.com/rc452860/vnet/utils"
)

type Record struct {
	ProxyAddr net.Addr
	ConnectionsPair
	Protocol   string
	Network    string
	RecordType string
	Up         uint64
	Down       uint64
	Status     int
}

type ConnectionsPair struct {
	Client net.Addr
	Target net.Addr
}

type GlobalResourceMonitor struct {
	sync.RWMutex
	LastMinuteOnline map[string][]ConnectionsPair
	TrafficUp        map[string]uint64
	TrafficDown      map[string]uint64
	TrafficUpSpeed   map[string]uint64
	TrafficDownSpeed map[string]uint64
	GlobalUp         uint64
	GlobalDown       uint64
	GlobalSpeedUp    uint64
	GlobalSpeeddown  uint64
	PacketCount      uint64
	// cancel           context.CancelFunc
	// context          context.Context
}

var (
	globalResourceMonitor *GlobalResourceMonitor
)

func GetGRMInstance() *GlobalResourceMonitor {
	utils.Lock("GRM")
	defer utils.UnLock("GRM")

	if globalResourceMonitor == nil {
		globalResourceMonitor = &GlobalResourceMonitor{
			LastMinuteOnline: make(map[string][]ConnectionsPair),
			TrafficUp:        make(map[string]uint64),
			TrafficDown:      make(map[string]uint64),
			TrafficUpSpeed:   make(map[string]uint64),
			TrafficDownSpeed: make(map[string]uint64),
			GlobalUp:         0,
			GlobalDown:       0,
			GlobalSpeedUp:    0,
			GlobalSpeeddown:  0,
			PacketCount:      0,
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	// globalResourceMonitor.cancel = cancel
	// globalResourceMonitor.context = ctx
	eventbus.GetEventBus().SubscribeAsync(
		"traffic",
	)
	return globalResourceMonitor
}

func (g *GlobalResourceMonitor) Stop() {
	utils.Lock("GRM")
	defer utils.UnLock("GRM")
	// g.cancel()
	globalResourceMonitor = nil

}
func (g *GlobalResourceMonitor) Subscript(name string, sub <-chan Record) {
	g.recordObservers.Store(name, sub)
}

func (g *GlobalResourceMonitor) UnSubscript(name string) {
	g.recordObservers.Delete(name)
}

func (g *GlobalResourceMonitor) Publish(record Record) {
	g.reciver <- record
}

func (g *GlobalResourceMonitor) process() {
	for {
		// var data interface{}
		// select {
		// case <-g.context.Done():
		// 	return
		// default:
		// case data = <-g.reciver:
		// }

	}
}

func (g *GlobalResourceMonitor) trafficStatistics(data Record) {
	g.TrafficDown[data.Port] += data.Down
	g.TrafficUp[data.Port] += data.Up
	g.GlobalUp += data.Up
	g.GlobalDown += data.Down

}
