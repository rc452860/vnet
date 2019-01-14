package proxy

import (
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
	Tcp          net.Listener        `json:"tcp"`
	Udp          net.PacketConn      `json:"udp"`
	UpSpeed      uint64              `json:"upspeed"`
	DownSpeed    uint64              `json:"downspeed"`
	UpBytes      uint64              `json:"upbytes"`
	DownBytes    uint64              `json:"downbytes"`
	TcpClose     chan struct{}       `json:"_"`
	UdpClose     chan struct{}       `json:"_"`
	TrafficClose chan struct{}       `json:"_"`
	TrafficMQ    chan TrafficMessage `json:"_"`
	TcpLock      *sync.Mutex         `json:"_"`
	UdpLock      *sync.Mutex         `json:"_"`
	Wait         *sync.WaitGroup     `json:"-"`
	Status       string              `json:"status"`
}

func NewProxyService() *ProxyService {
	return &ProxyService{
		TcpClose:     make(chan struct{}),
		UdpClose:     make(chan struct{}),
		TrafficClose: make(chan struct{}),
		TcpLock:      &sync.Mutex{},
		UdpLock:      &sync.Mutex{},
		TrafficMQ:    make(chan TrafficMessage, 128),
		Wait:         &sync.WaitGroup{},
		Status:       "stop",
	}
}

var trafficHandle func(TrafficMessage)

func RegisterTrafficHandle(trafficHandle func(TrafficMessage)) {
	trafficHandle = trafficHandle
}

func (this *ProxyService) TrafficMeasure() {
	speedClose := make(chan struct{})
	countClose := make(chan struct{})
	this.Wait.Add(1)
	go func() {
		var upTmp, downTmp uint64 = this.UpBytes, this.DownBytes
		tick := time.Tick(1 * time.Second)
		for {
			this.UpSpeed, upTmp = this.UpBytes-upTmp, this.UpBytes
			this.DownSpeed, downTmp = this.DownBytes-downTmp, this.DownBytes
			select {
			case <-tick:
				continue
			case <-speedClose:
				this.Wait.Done()
				return
			}
		}
	}()
	this.Wait.Add(1)
	go func() {
		for {
			var data TrafficMessage
			select {
			case <-countClose:
				log.Info("close countClose")
				this.Wait.Done()
				return
			case data = <-this.TrafficMQ:
			}
			this.UpBytes += data.UpBytes
			this.DownBytes += data.DownBytes
			if trafficHandle != nil {
				trafficHandle(data)
			}
		}
	}()
	<-this.TrafficClose
	log.Info("close traffic measure")
	speedClose <- struct{}{}
	countClose <- struct{}{}

}

func (this *ProxyService) Start() error {
	this.Status = "run"
	this.Wait.Add(2)
	return nil
}

func (this *ProxyService) Stop() error {
	log.Info("proxy stop")
	this.TrafficClose <- struct{}{}
	// this.TcpClose <- struct{}{}
	if this.Tcp != nil {
		err := this.Tcp.Close()
		if err != nil {
			log.Err(err)
		}
		this.Wait.Done()
	}
	if this.Udp != nil {
		err := this.Udp.Close()
		if err != nil {
			log.Err(err)
		}
		this.Wait.Done()
	}
	this.Wait.Wait()
	this.Status = "stop"
	// this.UdpClose <- struct{}{}
	return nil
}
