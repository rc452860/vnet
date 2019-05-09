package server

import (
	"github.com/rc452860/vnet/network"
	"github.com/rc452860/vnet/socks"
	"github.com/rc452860/vnet/utils/netx"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
	"net"
	"time"
)


// ShadowsocksProxy is respect shadowsocks proxy server
// it have Start and Stop method to control proxy
type ShadowsocksRProxy struct {
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	Method   string `json:"method,omitempty"`
	Password string `json:"password,omitempty"`
	ShadowsocksRArgs
	ReadLimiter   *rate.Limiter `json:"read_limit,omitempty"`
	WriteLimiter  *rate.Limiter `json:"write_limit,omitempty"`
	Protocol      string        `json:"protocol,omitempty"`
	ProtocolParam string        `json:"protocolParam,omitempty"`
	Obfs          string        `json:"obfs,omitempty"`
	ObfsParam     string        `json:"obfsParam,omitempty"`
	*network.Listener
}

// ShadowsocksArgs is ShadowsocksProxy arguments
type ShadowsocksRArgs struct {
	ConnectTimeout time.Duration `json:"connect_timeout,omitempty"`
	Limit          uint64        `json:"limit"`
	TCPSwitch      string        `json:"tcp_switch"`
	UDPSwitch      string        `json:"udp_switch"`
}

func (ssr *ShadowsocksRProxy) StartTCP() error {
	return ssr.ListenTCP(func(request *network.Request) {
		ssrd, err := network.NewShadowsocksRDecorate(request,
			ssr.Obfs, ssr.Method,
			ssr.Password, ssr.Protocol,
			ssr.ObfsParam, ssr.ProtocolParam,
			ssr.Host, ssr.Port,
			false)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"requestId": request.RequestID,
			}).Error("shadowsocksr NewShadowsocksRDecorate error")
		}
		go func() {
			defer func(){
				logrus.WithFields(logrus.Fields{
					"requestId": request.RequestID,
				}).Errorf("shadowsocksr connection read error %s",err)
			}()
			addr,err := socks.ReadAddr(ssrd)
			if err !=nil{
				logrus.WithFields(logrus.Fields{
					"requestId": ssrd.RequestID,
				}).Errorf("shadowsocksr read address error %s",err)
			}

			req,err := network.DialTcp(addr.String())
			if err != nil{
				logrus.WithFields(logrus.Fields{
					"requestId": ssrd.RequestID,
				}).Errorf("shadowsocksr proxy remote error %s",err)
			}
			defer req.Close()
			defer ssrd.Close()
			_ =req.SetKeepAlive(true)
			_,_,err = netx.DuplexCopyTcp(ssrd,req)
			if err != nil {
				if err, ok := err.(net.Error); ok && err.Timeout() {
					return // ignore i/o timeout
				}
				logrus.WithFields(logrus.Fields{
					"requestId": ssrd.RequestID,
				}).Errorf("shadowsocksr proxy process error %s",err)
			}
		}()

	})
}

func (ssr *ShadowsocksRProxy) StartUDP() {

}

func (ssr *ShadowsocksRProxy) TCPHandle(req *network.Request) {

}

func (ssr *ShadowsocksRProxy) handleStageAddr() {

}

func (ssr *ShadowsocksRProxy) handleStageConnection() {

}

func (ssr *ShadowsocksRProxy) Close() error{
	if ssr.TCP != nil{
		err := ssr.TCP.Close()
		if err !=nil{
			return err
		}
	}
	return nil
}