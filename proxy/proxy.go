package proxy

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/rc452860/vnet/log"
)

type IProxyService interface {
	Start() error
	Stop() error
}

type TrafficMessage struct {
	Network   string
	LAddr     string
	RAddr     string
	UpBytes   uint64
	DownBytes uint64
}
type ProxyService struct {
	context.Context `json:"-"`
	Tcp             net.Listener        `json:"tcp"`
	Udp             net.PacketConn      `json:"udp"`
	UpSpeed         uint64              `json:"upspeed"`
	DownSpeed       uint64              `json:"downspeed"`
	UpBytes         uint64              `json:"upbytes"`
	DownBytes       uint64              `json:"downbytes"`
	TrafficMQ       chan TrafficMessage `json:"_"`
	TcpLock         *sync.Mutex         `json:"_"`
	UdpLock         *sync.Mutex         `json:"_"`
	Status          string              `json:"status"`
	Cancel          context.CancelFunc  `json:"-`
}

func NewProxyService() *ProxyService {
	ctx, cancel := context.WithCancel(context.Background())
	return &ProxyService{
		TcpLock:   &sync.Mutex{},
		UdpLock:   &sync.Mutex{},
		TrafficMQ: make(chan TrafficMessage, 128),
		Status:    "stop",
		Context:   ctx,
		Cancel:    cancel,
	}
}

var trafficMonitorQueue []chan TrafficMessage

// RegisterTrafficHandle TrafficMessage channel
func RegisterTrafficHandle(trafficMonitor chan TrafficMessage) {
	trafficMonitorQueue = append(trafficMonitorQueue, trafficMonitor)
}

func (this *ProxyService) TrafficMeasure() {

	go func() {
		var upTmp, downTmp uint64 = this.UpBytes, this.DownBytes
		tick := time.Tick(1 * time.Second)
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
	}()
	go func() {
		for {
			var data TrafficMessage
			select {
			case <-this.Done():
				log.Info("close countClose")
				return
			case data = <-this.TrafficMQ:
			}
			this.UpBytes += data.UpBytes
			this.DownBytes += data.DownBytes
			if trafficMonitorQueue != nil && len(trafficMonitorQueue) > 0 {
				for _, item := range trafficMonitorQueue {
					item <- data
				}
			}
		}
	}()
	<-this.Done()
	log.Info("close traffic measure")

}

func (this *ProxyService) Start() error {
	this.Status = "run"
	return nil
}

func (this *ProxyService) Stop() error {
	log.Info("proxy stop")
	this.Cancel()
	// this.TcpClose <- struct{}{}
	if this.Tcp != nil {
		err := this.Tcp.Close()
		if err != nil {
			log.Err(err)
		}
	}
	if this.Udp != nil {
		err := this.Udp.Close()
		if err != nil {
			log.Err(err)
		}
	}
	this.Status = "stop"
	return nil
}
