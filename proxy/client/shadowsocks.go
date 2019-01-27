package client

import (
	"fmt"
	"net"

	"github.com/rc452860/vnet/common/log"
	"github.com/rc452860/vnet/proxy/common"

	"github.com/rc452860/vnet/socks"

	"github.com/rc452860/vnet/ciphers"

	"github.com/rc452860/vnet/conn"
)

type ShadowsocksClient struct {
	Host     string
	Method   string
	Password string
	Port     int
}

func NewShadowsocksClient(host, method, password string, port int) *ShadowsocksClient {
	return &ShadowsocksClient{
		Host:     host,
		Method:   method,
		Password: password,
		Port:     port,
	}
}

func (this *ShadowsocksClient) TcpProxy(con conn.IConn, target string, port int) error {
	proxy, err := net.Dial("tcp", fmt.Sprintf("%s:%v", this.Host, this.Port))
	if err != nil {
		log.Err(err)
		return err
	}
	proxys, err := conn.DefaultDecorate(proxy, conn.TCP)
	if err != nil {
		log.Err(err)
		return err
	}
	proxys, err = ciphers.CipherDecorate(this.Password, this.Method, proxys)
	if err != nil {
		log.Err(err)
		return err
	}
	addr := socks.ParseAddr(fmt.Sprintf("%s:%v", target, port))
	_, err = proxys.Write(addr)
	if err != nil {
		log.Err(err)
		return err
	}
	(&common.TcpChannel{}).Transport(con, proxys)
	return nil
}
