package server

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/rc452860/vnet/utils/addr"

	"github.com/rc452860/vnet/common/cache"
	"github.com/rc452860/vnet/common/eventbus"

	"github.com/rc452860/vnet/record"

	"github.com/rc452860/vnet/common/log"
)

type IProxyService interface {
	Start() error
	Stop() error
}

type ProxyService struct {
	context.Context          `json:"-"`
	TCP                      net.Listener       `json:"tcp"`
	UDP                      net.PacketConn     `json:"udp"`
	LastOneMinuteConnections *cache.Cache       `json:"-"`
	UpSpeed                  uint64             `json:"upspeed"`
	DownSpeed                uint64             `json:"downspeed"`
	UpBytes                  uint64             `json:"upbytes"`
	DownBytes                uint64             `json:"downbytes"`
	MessageRoute             chan interface{}   `json:"-"`
	Status                   string             `json:"status"`
	Cancel                   context.CancelFunc `json:"-"`
	Tick                     time.Duration      `json:"-"`
}

func NewProxyService() *ProxyService {
	return NewProxyServiceWithTick(time.Second)
}

func NewProxyServiceWithTick(duration time.Duration) *ProxyService {
	ctx, cancel := context.WithCancel(context.Background())
	return &ProxyService{
		MessageRoute:             make(chan interface{}, 32),
		Status:                   "stop",
		Context:                  ctx,
		Cancel:                   cancel,
		Tick:                     duration,
		LastOneMinuteConnections: cache.New(duration),
	}
}

var trafficMonitorQueue []chan record.Traffic

// RegisterTrafficHandle TrafficMessage channel
func RegisterTrafficHandle(trafficMonitor chan record.Traffic) {
	trafficMonitorQueue = append(trafficMonitorQueue, trafficMonitor)
}

func (this *ProxyService) TrafficMeasure() {
	go this.speed()
	go this.route()
	<-this.Done()
	log.Info("close traffic measure")
}

func (this *ProxyService) route() {
	for {
		var data interface{}
		select {
		case <-this.Done():
			log.Info("close countClose")
			return
		case data = <-this.MessageRoute:
		}
		// handle data is null cause type switch fail.
		if data == nil {
			continue
		}
		switch data.(type) {
		case record.Traffic:
			this.traffic(data.(record.Traffic))
		case record.ConnectionProxyRequest:
			this.proxyRequest(data.(record.ConnectionProxyRequest))
		}
	}
}

// traffic handle traffic message
func (this *ProxyService) traffic(data record.Traffic) {
	eventbus.GetEventBus().Publish("record:traffic", data)
	this.UpBytes += data.Up
	this.DownBytes += data.Down
	if trafficMonitorQueue != nil && len(trafficMonitorQueue) > 0 {
		for _, item := range trafficMonitorQueue {
			item <- data
		}
	}
}

// proxyRequest handle proxy request message
func (this *ProxyService) proxyRequest(data record.ConnectionProxyRequest) {
	eventbus.GetEventBus().Publish("record:proxyRequest", data)
	key := addr.GetIPFromAddr(data.ClientAddr)
	if this.LastOneMinuteConnections.Get(key) == nil {
		this.LastOneMinuteConnections.Put(key, []record.ConnectionProxyRequest{data}, this.Tick)
	} else {
		last := this.LastOneMinuteConnections.Get(key).([]record.ConnectionProxyRequest)
		this.LastOneMinuteConnections.Put(key, append(last, data), this.Tick)
	}

	// just print tcp log
	if strings.Contains("tcp", data.ClientAddr.Network()) {
		log.Info("%s <----%v:%s----> %s",
			data.ClientAddr.String(),
			addr.GetPortFromAddr(data.ProxyAddr),
			addr.GetNetworkFromAddr(data.ProxyAddr),
			fmt.Sprintf("%s:%v", data.GetAddress(), data.GetPort()))
	}

}

// speed is traffic statis
func (this *ProxyService) speed() {
	var upTmp, downTmp uint64 = this.UpBytes, this.DownBytes
	tick := time.Tick(this.Tick)
	for {
		this.UpSpeed, upTmp = this.UpBytes-upTmp, this.UpBytes
		this.DownSpeed, downTmp = this.DownBytes-downTmp, this.DownBytes
		select {
		case <-tick:
			continue
		case <-this.Done():
			return
		}
	}
}

func (this *ProxyService) Start() error {
	go this.TrafficMeasure()
	this.Status = "run"
	return nil
}

func (this *ProxyService) Stop() error {
	log.Info("proxy stop")
	this.Cancel()
	if this.TCP != nil {
		err := this.TCP.Close()
		if err != nil {
			log.Err(err)
		}
	}
	if this.UDP != nil {
		err := this.UDP.Close()
		if err != nil {
			log.Err(err)
		}
	}
	this.Status = "stop"
	return nil
}
