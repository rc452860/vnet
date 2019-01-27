package record

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/rc452860/vnet/common/cache"
	"github.com/rc452860/vnet/common/eventbus"
	"github.com/rc452860/vnet/common/log"
	"github.com/rc452860/vnet/utils"
	"github.com/rc452860/vnet/utils/addr"
)

// Traffic is repersent traffic record
type Traffic struct {
	ConnectionPair `json:"connections_pair,omitempty"`
	Network        string `json:"network,omitempty"`
	Up             uint64 `json:"up,omitempty"`
	Down           uint64 `json:"down,omitempty"`
}

type ConnectionProxyRequest struct {
	ConnectionPair
}

// ConnectionPair is struct represent a proxy connections
type ConnectionPair struct {
	ProxyAddr  net.Addr `json:"proxy_addr,omitempty"`
	ClientAddr net.Addr `json:"client_addr,omitempty"`
	TargetAddr net.Addr `json:"target_addr,omitempty"`
}

// GlobalResourceMonitor global resource monitor
type GlobalResourceMonitor struct {
	sync.RWMutex
	LastOneMinuteConnections *cache.Cache `json:"-"`
	GlobalUp                 uint64       `json:"global_up,omitempty"`
	GlobalDown               uint64       `json:"global_down,omitempty"`
	GlobalSpeedUp            uint64       `json:"global_speed_up,omitempty"`
	GlobalSpeeddown          uint64       `json:"global_speeddown,omitempty"`
	PacketCount              uint64       `json:"packet_count,omitempty"`
	tick                     time.Duration
}

var (
	globalResourceMonitor *GlobalResourceMonitor
)

// GetGRMInstance GlobalResourceMonitor singal instance
func GetGRMInstance() *GlobalResourceMonitor {
	return GetGRMInstanceWithTick(time.Minute)
}

func GetGRMInstanceWithTick(duration time.Duration) *GlobalResourceMonitor {
	if globalResourceMonitor != nil {
		return globalResourceMonitor
	}
	utils.Lock("GetGRMInstance")
	globalResourceMonitor = &GlobalResourceMonitor{
		LastOneMinuteConnections: cache.New(duration),
		GlobalUp:                 0,
		GlobalDown:               0,
		GlobalSpeedUp:            0,
		GlobalSpeeddown:          0,
		PacketCount:              0,
		tick:                     duration,
	}
	eventbus.GetEventBus().SubscribeAsync("record:traffic", globalResourceMonitor.trafficStatistics, false)
	eventbus.GetEventBus().SubscribeAsync("record:proxyRequest", globalResourceMonitor.lastOneMinuteConnections, false)
	go globalResourceMonitor.speed()
	utils.UnLock("GetGRMInstance")
	return globalResourceMonitor
}

func (g *GlobalResourceMonitor) trafficStatistics(data Traffic) {
	utils.Lock("record:globalResourceMonitor")
	defer utils.UnLock("record:globalResourceMonitor")
	g.GlobalUp += data.Up
	g.GlobalDown += data.Down
	if strings.Contains(data.Network, "udp") {
		g.PacketCount++
	}
}
func (g *GlobalResourceMonitor) lastOneMinuteConnections(data ConnectionProxyRequest) {
	key := addr.GetIPFromAddr(data.ClientAddr)
	if g.LastOneMinuteConnections.Get(key) == nil {
		g.LastOneMinuteConnections.Put(key, []ConnectionProxyRequest{data}, g.tick)
	} else {
		last := g.LastOneMinuteConnections.Get(key).([]ConnectionProxyRequest)
		g.LastOneMinuteConnections.Put(key, append(last, data), g.tick)
	}
}

// speed is initialize by GlobalResourceMonitor init,and it should be in other goroutine
// it's a schedule for calculate up and down speed for global and every single service
func (g *GlobalResourceMonitor) speed() {
	var globalUpTmp, globalDownTmp uint64
	tick := time.Tick(g.tick)
	for {
		<-tick
		utils.Lock("record:globalResourceMonitor")
		g.GlobalSpeedUp = g.GlobalUp - globalUpTmp
		globalUpTmp = g.GlobalUp
		g.GlobalSpeeddown = g.GlobalDown - globalDownTmp
		globalDownTmp = g.GlobalDown
		utils.UnLock("record:globalResourceMonitor")
	}
}

// GetLastOneMinuteOnlineCount return online service count
func (g *GlobalResourceMonitor) GetLastOneMinuteOnlineCount() int {
	return g.LastOneMinuteConnections.Size()
}

// GetLastOneMinuteOnlineByPort convert LastOneMinuteOnlineByPort key from client to proxy
// the original map is map[client.ip]
// the convert map is map[proxy.port]
func (g *GlobalResourceMonitor) GetLastOneMinuteOnlineByPort() map[int][]net.Addr {
	result := make(map[int][]net.Addr)
	filter := make(map[string]bool)
	g.LastOneMinuteConnections.Range(func(k, v interface{}) {
		if value, ok := v.([]ConnectionProxyRequest); ok {
			for _, item := range value {
				proxyPort := addr.GetPortFromAddr(item.ProxyAddr)
				clientIP := addr.GetIPFromAddr(item.ClientAddr)
				filterKey := fmt.Sprintf("%v-%s", proxyPort, clientIP)
				if filter[filterKey] {
					continue
				} else {
					filter[filterKey] = true
				}
				if result[proxyPort] != nil {
					result[proxyPort] = append(result[proxyPort],
						item.ClientAddr)

				} else {
					result[proxyPort] = []net.Addr{item.ClientAddr}
				}
			}
		}
	})
	return result
}

func (g *GlobalResourceMonitor) String() string {
	result, e := json.Marshal(g)
	if e != nil {
		log.Err(e)
		return ""
	}
	return string(result)
}
