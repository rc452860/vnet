package proxy

import (
	"net"
	"sync"
)

type IProxyService interface {
	Start() error
	Stop() error
}

type ProxyService struct {
	Tcp      net.Listener
	Udp      net.PacketConn
	TcpClose chan struct{}
	UdpClose chan struct{}
	tcpLock  *sync.Mutex
	udpLock  *sync.Mutex
}

func NewProxyService() ProxyService {
	return ProxyService{
		TcpClose: make(chan struct{}),
		UdpClose: make(chan struct{}),
		tcpLock:  &sync.Mutex{},
		udpLock:  &sync.Mutex{},
	}
}

func (this ProxyService) Start() error {
	return nil
}

func (this ProxyService) Stop() error {
	this.UdpClose <- struct{}{}
	this.UdpClose <- struct{}{}
	return nil
}
