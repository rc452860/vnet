package proxy

import (
	"encoding/hex"
	"github.com/pkg/errors"
	"github.com/rc452860/vnet/common"
	"github.com/rc452860/vnet/common/log"
	"github.com/rc452860/vnet/common/network"
	"github.com/rc452860/vnet/common/pool"
	"github.com/rc452860/vnet/utils/binaryx"
	"github.com/rc452860/vnet/utils/goroutine"
	"github.com/rc452860/vnet/utils/netx"
	"github.com/rc452860/vnet/utils/socksproxy"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
	"net"
	"runtime/debug"
	"sync"
	"time"
)

// ShadowsocksProxy is respect shadowsocks proxy service
// it have Start and Stop method to control proxy
type ShadowsocksRProxy struct {
	Host              string `json:"host,omitempty"`
	Port              int    `json:"port,omitempty"`
	Method            string `json:"method,omitempty"`
	Password          string `json:"password,omitempty"`
	Protocol          string `json:"protocol,omitempty"`
	ProtocolParam     string `json:"protocolParam,omitempty"`
	Obfs              string `json:"obfs,omitempty"`
	ObfsParam         string `json:"obfsParam,omitempty"`
	*network.Listener `json:"-"`
	Users             map[string]string      `json:"users,omitempty"`
	Status            string                 `json:"status,omitempty"`
	common.TrafficReport `json:"-"`
	common.OnlineReport `json:"-"`
	*ShadowsocksRArgs
}

// ShadowsocksArgs is ShadowsocksProxy arguments
type ShadowsocksRArgs struct {
	ConnectTimeout time.Duration `json:"connect_timeout,omitempty"`
	Limiter        *rate.Limiter
	TCPSwitch      string `json:"tcp_switch"`
	UDPSwitch      string `json:"udp_switch"`
}

// Start tcp and udp according to the configuration
func (ssr *ShadowsocksRProxy) Start() error {
	var err error
	if ssr.ShadowsocksRArgs.TCPSwitch != "false" {
		err = ssr.StartTCP()
		if err != nil {
			return err
		}
	}

	if ssr.ShadowsocksRArgs.UDPSwitch != "false" {
		err = ssr.StartUDP()
		if err != nil {
			return err
		}
	}
	return nil
}

func (ssr *ShadowsocksRProxy) Stop() error{
	return ssr.Listener.Close()
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
			return
		}
		ssrd.TrafficReport = ssr.TrafficReport
		go func() {
			defer func() {
				if err := recover(); err != nil {
					logrus.WithFields(logrus.Fields{
						"requestId": request.RequestID,
					}).Errorf("shadowsocksr connection read error :%v stack: %s", err, string(debug.Stack()))
				}
			}()
			addr, err := socksproxy.ReadAddr(ssrd)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"requestId": ssrd.RequestID,
				}).Errorf("shadowsocksr read address error %s", err)
				return
			}
			ssr.handleStageAddr(ssrd.UID,ssrd.RemoteAddr().String(), ssrd.LocalAddr().String(), addr.String(), "tcp")
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

func (ssr *ShadowsocksRProxy) StartUDP() error {
	err := ssr.ListenUDP(func(request *network.Request) {
		go func() {
			defer func() {
				if e := recover(); e != nil {
					logrus.Errorf("shadowsocksr udp listener crashed , err : %s , \ntrace:%s", e, string(debug.Stack()))
				}
			}()
			ssrd, err := network.NewShadowsocksRDecorate(request,
				ssr.Obfs, ssr.Method,
				ssr.Password, ssr.Protocol,
				ssr.ObfsParam, ssr.ProtocolParam,
				ssr.Host, ssr.Port,
				false,
				ssr.Users)
			ssrd.TrafficReport =  ssr.TrafficReport
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"requestId": request.RequestID,
					"error":     err,
				}).Error("shadowsocksr NewShadowsocksRDecorate error")
			}
			// TODO UDP TIMEOUT
			udpMap := NewShadowsocksRUDPMap(30)
			for {
				data, uid, addr, err := ssrd.ReadFrom()
				if err != nil {
					logrus.WithFields(logrus.Fields{
						"err": err,
					}).Error("ShadowsocksRDecrate read udp error")
					continue
				}
				remoteAddr := socksproxy.SplitAddr(data)
				if remoteAddr == nil {
					continue
				}
				logrus.WithFields(logrus.Fields{
					"remoteAddr": remoteAddr.String(),
					"serverAddr": ssrd.PacketConn.LocalAddr().String(),
					"clientAddr": addr.String(),
					"uid":        binaryx.LEBytesToUInt32(uid),
				}).Info("recive udp proxy")
				data = data[len(remoteAddr.Raw):]
				remotePacketConn := udpMap.Get(addr.String())
				if remotePacketConn == nil {
					remotePacketConn = &ShadowsocksRUDPMapItem{}
					remotePacketConn.Uid = uid
					remotePacketConn.PacketConn, err = net.ListenPacket("udp", "")
					if err != nil {
						logrus.WithFields(logrus.Fields{
							"remoteAddr": remoteAddr.String(),
							"serverAddr": ssrd.PacketConn.LocalAddr().String(),
							"clientAddr": addr.String(),
							"uid":        binaryx.LEBytesToUInt32(uid),
							"err":        err,
						}).Error("shadowoscksr listenPacket udp error")
						continue
					}
					udpMap.Add(addr, ssrd, remotePacketConn)
				}
				remoteAddrResolve, err := net.ResolveUDPAddr("udp", remoteAddr.String())
				if err != nil {
					logrus.WithFields(logrus.Fields{
						"remoteAddr": remoteAddr.String(),
						"serverAddr": ssrd.PacketConn.LocalAddr().String(),
						"clientAddr": addr.String(),
						"uid":        binaryx.LEBytesToUInt32(uid),
						"err":        err,
					}).Error("shadowoscksr listenPacket udp error")
					continue
				}
				ssr.handleStageAddr(int(binaryx.LEBytesToUInt32(uid)),addr.String(), ssrd.PacketConn.LocalAddr().String(), remoteAddr.String(), "udp")
				udpMap.Add(addr, ssrd, remotePacketConn)
				_, err = remotePacketConn.WriteTo(data, remoteAddrResolve)
				if err != nil {
					if err != nil {
						logrus.WithFields(logrus.Fields{
							"remoteAddr": remoteAddr.String(),
							"serverAddr": ssrd.PacketConn.LocalAddr().String(),
							"clientAddr": addr.String(),
							"uid":        binaryx.LEBytesToUInt32(uid),
							"err":        err,
						}).Error("shadowoscksr listenPacket udp error")
						continue
					}
					udpMap.Add(addr, ssrd, remotePacketConn)
				}
			}
		}()
	})
	return err
}

func (ssr *ShadowsocksRProxy) handleStageAddr(uid int,client, server, proxyTarget, network string) {
	if uid == 0{
		logrus.WithFields(logrus.Fields{
			"uid":uid,
			"client":client,
			"server":server,
			"proxyTarget":proxyTarget,
		}).Warn("handleStageAddr uid is 0")
		return
	}
	ssr.OnlineReport.Online(uid,client)
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

func (ssr *ShadowsocksRProxy) AddUser(uid int, password string) {
	if ssr.Users == nil {
		ssr.Users = make(map[string]string)
	}
	uidPack := binaryx.LEUint32ToBytes(uint32(uid))
	logrus.Debugf("shadowsocksr adduser uidPack: %s", hex.EncodeToString(uidPack))
	uidPackStr := string(uidPack)
	ssr.Users[uidPackStr] = password
}

func (ssr *ShadowsocksRProxy) DelUser(uid int) {
	if ssr.Users == nil {
		return
	}
	uidPack := string(binaryx.LEUint32ToBytes(uint32(uid)))
	delete(ssr.Users, uidPack)
}

func (ssr *ShadowsocksRProxy) Reload(users map[string]string) {
	ssr.Users = users
}

type ShadowsocksRUDPMapItem struct {
	net.PacketConn
	Uid []byte
}

// Packet NAT table
type ShadowsocksRUDPMap struct {
	sync.RWMutex
	m       map[string]*ShadowsocksRUDPMapItem
	timeout time.Duration
}

func NewShadowsocksRUDPMap(timeout time.Duration) *ShadowsocksRUDPMap {
	m := &ShadowsocksRUDPMap{}
	m.m = make(map[string]*ShadowsocksRUDPMapItem)
	m.timeout = timeout
	return m
}

func (m *ShadowsocksRUDPMap) Get(key string) *ShadowsocksRUDPMapItem {
	m.RLock()
	defer m.RUnlock()
	return m.m[key]
}

func (m *ShadowsocksRUDPMap) Set(key string, pc *ShadowsocksRUDPMapItem) {
	m.Lock()
	defer m.Unlock()
	m.m[key] = pc
}

func (m *ShadowsocksRUDPMap) Del(key string) *ShadowsocksRUDPMapItem {
	m.Lock()
	defer m.Unlock()

	pc, ok := m.m[key]
	if ok {
		delete(m.m, key)
		return pc
	}
	return nil
}

func (m *ShadowsocksRUDPMap) Add(client net.Addr, server *network.ShadowsocksRDecorate, remoteServer *ShadowsocksRUDPMapItem) {
	m.Set(client.String(), remoteServer)
	go goroutine.Protect(func() {
		//TODO defer recover
		_ = ShadowsocksRMapTimeCopy(server, client, remoteServer, m.timeout)
		if pc := m.Del(client.String()); pc != nil {
			_ = pc.Close()
		}
	})
}

// copy from src to dst at target with read timeout
func ShadowsocksRMapTimeCopy(dst *network.ShadowsocksRDecorate, target net.Addr, src *ShadowsocksRUDPMapItem, timeout time.Duration) error {
	buf := pool.GetBuf()
	defer pool.PutBuf(buf)
	defer func() {
		if e := recover(); e != nil {
			log.Error("panic in timedCopy: %v", e)
		}
	}()

	for {
		_ = src.SetReadDeadline(time.Now().Add(timeout * time.Second))
		n, raddr, err := src.ReadFrom(buf)
		if err != nil {
			return errors.Cause(err)
		}

		srcAddr := socksproxy.ParseAddr(raddr.String())
		srcAddrByte := srcAddr.Raw
		copy(buf[len(srcAddrByte):], buf[:n])
		copy(buf, srcAddrByte)
		err = dst.WriteTo(buf[:len(srcAddrByte)+n], src.Uid, target)

		if err != nil {
			return errors.Cause(err)
		}
	}
}
