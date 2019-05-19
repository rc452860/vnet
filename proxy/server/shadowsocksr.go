package server

import (
	"encoding/hex"
	"github.com/rc452860/vnet/network"
	"github.com/rc452860/vnet/socks"
	"github.com/rc452860/vnet/utils/binaryx"
	"github.com/rc452860/vnet/utils/netx"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
	"net"
	"runtime/debug"
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
	ReadLimiter       *rate.Limiter `json:"read_limit,omitempty"`
	WriteLimiter      *rate.Limiter `json:"write_limit,omitempty"`
	Protocol          string        `json:"protocol,omitempty"`
	ProtocolParam     string        `json:"protocolParam,omitempty"`
	Obfs              string        `json:"obfs,omitempty"`
	ObfsParam         string        `json:"obfsParam,omitempty"`
	*network.Listener `json:"_"`
	Users             map[string]string `json:"users,omitempty"`
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
			false,
			ssr.Users)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"requestId": request.RequestID,
				"error":     err,
			}).Error("shadowsocksr NewShadowsocksRDecorate error")
		}
		go func() {
			defer func() {
				if err := recover(); err != nil {
					logrus.WithFields(logrus.Fields{
						"requestId": request.RequestID,
					}).Errorf("shadowsocksr connection read error :%v stack: %s", err, string(debug.Stack()))
				}
			}()
			addr, err := socks.ReadAddr(ssrd)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"requestId": ssrd.RequestID,
				}).Errorf("shadowsocksr read address error %s", err)
				return
			}
			logrus.Infof("reslove addr success: %s", addr.String())
			req, err := network.DialTcp(addr.String())
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"requestId": ssrd.RequestID,
				}).Errorf("shadowsocksr proxy remote error %s", err)
				return
			}
			defer req.Close()
			defer ssrd.Close()
			_ = req.SetKeepAlive(true)
			_, _, err = netx.DuplexCopyTcp(ssrd, req)
			if err != nil {
				if err, ok := err.(net.Error); ok && err.Timeout() {
					return // ignore i/o timeout
				}
				logrus.WithFields(logrus.Fields{
					"requestId": ssrd.RequestID,
				}).Errorf("shadowsocksr proxy process error %s", err)
				return
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

func (ssr *ShadowsocksRProxy) Close() error {
	if ssr.TCP != nil {
		err := ssr.TCP.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (ssr *ShadowsocksRProxy) AddUser(uid uint32, password string) {
	if ssr.Users == nil {
		ssr.Users = make(map[string]string)
	}
	uidPack := binaryx.LEUint32ToBytes(uid)
	logrus.Debugf("shadowsocksr adduser uidPack: %s", hex.EncodeToString(uidPack))
	uidPackStr := string(uidPack)
	ssr.Users[uidPackStr] = password
}

func (ssr *ShadowsocksRProxy) DelUser(uid uint32) {
	if ssr.Users == nil {
		return
	}
	uidPack := string(binaryx.LEUint32ToBytes(uid))
	delete(ssr.Users, uidPack)
}

func (ssr *ShadowsocksRProxy) Reload(users map[string]string) {
	ssr.Users = users
}
