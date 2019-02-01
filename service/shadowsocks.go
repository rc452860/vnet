package service

import (
	"fmt"

	"github.com/rc452860/vnet/proxy/server"
)

var ss *ShadowsocksService

func CurrentShadowsocksService() *ShadowsocksService {
	if ss == nil {
		ss = NewShadowsocksService()
	}
	return ss
}

type ShadowsocksService struct {
	Servers map[int]*server.ShadowsocksProxy
}

func NewShadowsocksService() *ShadowsocksService {
	return &ShadowsocksService{
		Servers: make(map[int]*server.ShadowsocksProxy),
	}
}

func (this *ShadowsocksService) Add(host string, method string, password string, port int, args server.ShadowsocksArgs) error {
	proxy := this.Get(port)

	if proxy == nil {
		proxy, err := server.NewShadowsocks(host, method, password, port, args)
		if err != nil {
			return err
		}
		this.Servers[port] = proxy
		return err
	} else {
		proxy.Host = host
		proxy.Method = method
		proxy.Password = password
		proxy.ShadowsocksArgs = args
		return nil
	}
}

func (this *ShadowsocksService) List() []*server.ShadowsocksProxy {
	list := make([]*server.ShadowsocksProxy, len(this.Servers))
	i := 0
	for _, v := range this.Servers {
		list[i] = v
		i++
	}
	return list
}

func (this *ShadowsocksService) Start(port int) error {
	ss := this.Servers[port]
	if ss != nil {
		ss.Start()
		return nil
	}
	return fmt.Errorf("port not found")
}

func (this *ShadowsocksService) Get(port int) *server.ShadowsocksProxy {
	return this.Servers[port]
}

func (this *ShadowsocksService) Stop(port int) error {
	server := this.Servers[port]
	if server != nil {
		return this.Servers[port].Stop()
	}
	return nil
}

func (this *ShadowsocksService) Del(port int) error {
	server := this.Servers[port]
	if server != nil && server.Status == "run" {
		err := server.Stop()
		if err != nil {
			return err
		}
	}
	delete(this.Servers, port)
	return nil
}

func (this *ShadowsocksService) IsExist(port int) bool {
	return this.Servers[port] == nil
}
