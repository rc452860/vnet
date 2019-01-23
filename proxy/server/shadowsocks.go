package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/rc452860/vnet/utils/datasize"

	"github.com/rc452860/vnet/ciphers"
	"github.com/rc452860/vnet/comm/log"
	"github.com/rc452860/vnet/conn"
	"github.com/rc452860/vnet/pool"
	"github.com/rc452860/vnet/proxy"
	"github.com/rc452860/vnet/record"
	"github.com/rc452860/vnet/socks"
	"golang.org/x/time/rate"
)

type mode int

const (
	remoteServer mode = iota
	relayClient
	socksClient
)

var (
	logging *log.Logging
	// ShadowsocksServerList Global map for store shadowsocks proxy
	ShadowsocksServerList map[string]*ShadowsocksProxy
)

func init() {
	logging = log.GetLogger("root")
}

// ShadowsocksProxy is respect shadowsocks proxy server
// it have Start and Stop method to control proxy
type ShadowsocksProxy struct {
	*proxy.ProxyService `json:"-,omitempty"`
	Host                string        `json:"host,omitempty"`
	Port                int           `json:"port,omitempty"`
	Method              string        `json:"method,omitempty"`
	Password            string        `json:"password,omitempty"`
	TCPTimeout          time.Duration `json:"tcp_timeout,omitempty"`
	UDPTimeout          time.Duration `json:"udp_timeout,omitempty"`
	ReadLimit           *rate.Limiter `json:"read_limit,omitempty"`
	WriteLimit          *rate.Limiter `json:"write_limit,omitempty"`
}

// NewShadowsocks is new ShadowsocksProxy object
func NewShadowsocks(host string, method string, password string, port int, limit string, timeout time.Duration) (*ShadowsocksProxy, error) {
	ss := &ShadowsocksProxy{
		ProxyService: proxy.NewProxyService(),
		Host:         host,
		Method:       method,
		Password:     password,
		Port:         port,
	}
	// for traffic limit
	err := ss.ConfigLimitHuman(limit)
	if err != nil {
		return nil, err
	}

	// for 3 second time out
	err = ss.ConfigTimeout(timeout)
	return ss, err
}

// ConfigLimit config shadowsocks traffic limit
func (s *ShadowsocksProxy) ConfigLimit(limit uint64) {
	if limit == 0 {
		return
	}
	s.ReadLimit = rate.NewLimiter(rate.Limit(limit), int(limit))
	s.WriteLimit = rate.NewLimiter(rate.Limit(limit), int(limit))
}

//ConfigLimitHuman is config shadowsocks service traffic limit
// argument limit can be a human readable like: 4KB 4MB 4GB
func (s *ShadowsocksProxy) ConfigLimitHuman(limit string) error {
	if limit != "" {
		trafficLimit, err := datasize.Parse(limit)
		if err != nil {
			log.Err(err)
			return err
		}
		if trafficLimit < 5*1024 {
			return nil
		}
		// logging.Info("server port: %v limit is: %v", s.Port, trafficLimit)
		s.ReadLimit = rate.NewLimiter(rate.Limit(trafficLimit), int(trafficLimit))
		s.WriteLimit = rate.NewLimiter(rate.Limit(trafficLimit), int(trafficLimit))
	}
	return nil

}

// ConfigTimeout is config shadowsocks timeout
func (s *ShadowsocksProxy) ConfigTimeout(timeout time.Duration) error {
	if timeout == 0 {
		s.TCPTimeout = 3e9
		s.UDPTimeout = 3e9
	} else {
		s.TCPTimeout = timeout
		s.UDPTimeout = timeout
	}
	// log.Info("%s:%v timeout:%v", s.Host, s.Port, s.TCPTimeout)
	return nil
}

// Start proxy
func (s *ShadowsocksProxy) Start() error {
	s.ProxyService.Start()
	if err := s.startTCP(); err != nil {
		log.Err(err)
		return err
	}
	if err := s.startUDP(); err != nil {
		log.Err(err)
		return err
	}
	return nil
}

// Stop proxy
func (s *ShadowsocksProxy) Stop() error {
	return s.ProxyService.Stop()
}

// statistics tcpUpload traffic
func (s *ShadowsocksProxy) tcpUpload(con conn.IConn, up uint64) {
	message := record.Traffic{
		ConnectionPair: record.ConnectionPair{
			ProxyAddr:  con.LocalAddr(),
			ClientAddr: con.RemoteAddr(),
		},
		Network: con.GetNetwork(),
		Up:      up,
	}

	s.ProxyService.MessageRoute <- message
}

// statics tcpDownload traffic
func (s *ShadowsocksProxy) tcpDownload(con conn.IConn, down uint64) {
	message := record.Traffic{
		ConnectionPair: record.ConnectionPair{
			ProxyAddr:  con.LocalAddr(),
			ClientAddr: con.RemoteAddr(),
		},
		Network: con.GetNetwork(),
		Down:    down,
	}

	s.ProxyService.MessageRoute <- message
}

// start shadowsocks tcp proxy service
func (s *ShadowsocksProxy) startTCP() error {
	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		logging.Error(err.Error())
		return err
	}
	server, err := net.ListenTCP("tcp", tcpAddr)
	// logging.Info("listening TCP on %s", addr)
	if err != nil {
		logging.Error(err.Error())
		return errors.Cause(err)
	}
	s.TCP = server

	go func() {
		defer server.Close()
		for {
			select {
			case <-s.ProxyService.Done():
				return
			default:
			}
			server.SetDeadline(time.Now().Add(s.TCPTimeout))
			lcon, err := server.Accept()
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}

			if err != nil {
				if strings.Contains(err.Error(), "use of closed network connection") {
					return
				}
				logging.Error(err.Error())
				continue
			}

			go func() {
				defer lcon.Close()
				/** 默认装饰器 */
				lcd, err := conn.DefaultDecorate(lcon, conn.TCP)
				if err != nil {
					logging.Err(err)
					return
				}
				/** 去皮流量记录装饰器 */
				lcd, err = conn.TrafficDecorate(lcd, s.tcpUpload, s.tcpDownload)
				if err != nil {
					logging.Err(err)
					return
				}
				/** 限流装饰器 */
				lcd, _ = conn.TrafficLimitDecorate(lcd, s.ReadLimit, s.WriteLimit)

				/** 加密装饰器 */
				lcd, err = ciphers.CipherDecorate(s.Password, s.Method, lcd)
				if err != nil {
					logging.Err(err)
					return
				}

				/** 读取目标地址 */
				targetAddr, err := socks.ReadAddr(lcd)
				if err != nil {
					log.Error("read target address error %s. (maybe the crypto method 2rong configuration)", err.Error())
					return
				}

				rc, err := net.Dial("tcp", targetAddr.String())
				if err != nil {
					logging.Error("connect target:%s error", targetAddr)
					logging.Err(err)
					return
				}
				defer rc.Close()
				s.ConnectionStage(s.TCP.Addr(), lcd.RemoteAddr(), rc.RemoteAddr())

				rc.(*net.TCPConn).SetKeepAlive(true)
				// logging.Info("tcp %s <----> %s", lcd.RemoteAddr(), targetAddr)

				/** 默认装饰器 */
				rcd, err := conn.DefaultDecorate(rc, conn.TCP)
				if err != nil {
					logging.Err(err)
					return
				}

				_, _, err = relayTCP(lcd, rcd)
				if err != nil {
					if err, ok := err.(net.Error); ok && err.Timeout() {
						return // ignore i/o timeout
					}
					logging.Error("relay error: %v", err)
				}
			}()
		}
	}()
	return nil
}

// relay copies between left and right bidirectionally. Returns number of
// bytes copied from right to left, from left to right, and any error occurred.
func relayTCP(left, right net.Conn) (int64, int64, error) {
	type res struct {
		N   int64
		Err error
	}
	ch := make(chan res)

	go func() {
		n, err := io.Copy(right, left)
		right.SetDeadline(time.Now()) // wake up the other goroutine blocking on right
		left.SetDeadline(time.Now())  // wake up the other goroutine blocking on left
		ch <- res{n, err}
	}()

	n, err := io.Copy(left, right)
	right.SetDeadline(time.Now()) // wake up the other goroutine blocking on right
	left.SetDeadline(time.Now())  // wake up the other goroutine blocking on left
	rs := <-ch

	if err == nil {
		err = rs.Err
	}
	return n, rs.N, errors.Cause(err)
}

// udp upload traffic count
func (s *ShadowsocksProxy) udpUpload(laddr, raddr net.Addr, n uint64) {
	message := record.Traffic{
		ConnectionPair: record.ConnectionPair{
			ProxyAddr:  laddr,
			ClientAddr: raddr,
		},
		Network: laddr.Network(),
		Up:      n,
	}

	s.ProxyService.MessageRoute <- message
}

// udp download traffic count
func (s *ShadowsocksProxy) udpDownload(laddr, raddr net.Addr, n uint64) {
	message := record.Traffic{
		ConnectionPair: record.ConnectionPair{
			ProxyAddr:  laddr,
			ClientAddr: raddr,
		},
		Network: laddr.Network(),
		Down:    n,
	}
	s.ProxyService.MessageRoute <- message
}

// Listen on addr for encrypted packets and basically do UDP NAT.
func (s *ShadowsocksProxy) startUDP() error {
	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)
	server, err := net.ListenPacket("udp", addr)
	if err != nil {
		logging.Error("UDP remote listen error: %v", err)
		return errors.Cause(err)
	}
	s.UDP = server
	// 去皮流量装饰器
	server = conn.PacketTrafficConnDecorate(server, s.udpUpload, s.udpDownload)
	server, err = ciphers.CipherPacketDecorate(s.Password, s.Method, server)
	if err != nil {
		logging.Error("UDP CipherPacketDecorate init error: %v", err)
		return errors.Cause(err)
	}

	nm := newNATmap(s.UDPTimeout)
	buf := pool.GetUdpBuf()
	defer pool.PutUdpBuf(buf)

	// logging.Info("listening UDP on %s", addr)

	go func() {
		defer server.Close()
		for {
			select {
			case <-s.ProxyService.Done():
				return
			default:
			}
			server.SetDeadline(time.Now().Add(s.UDPTimeout))
			n, raddr, err := server.ReadFrom(buf)
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			if err != nil {
				if strings.Contains(err.Error(), "use of closed network connection") {
					return
				}
				logging.Error("UDP remote read error: %v", err)
				continue
			}
			tgtAddr := socks.SplitAddr(buf[:n])
			if tgtAddr == nil {
				logging.Error("failed to split target address from packet: %q", buf[:n])
				continue
			}
			// logging.Info("udp %s <----> %s", raddr, tgtAddr)
			tgtUDPAddr, err := net.ResolveUDPAddr("udp", tgtAddr.String())
			if err != nil {
				logging.Error("failed to resolve target UDP address: %v", err)
				continue
			}

			s.ConnectionStage(s.UDP.LocalAddr(), raddr, tgtUDPAddr)
			payload := buf[len(tgtAddr):n]

			pc := nm.Get(raddr.String())
			if pc == nil {
				pc, err = net.ListenPacket("udp", "")
				if err != nil {
					logging.Error("UDP remote listen error: %v", err)
					continue
				}

				nm.Add(raddr, server, pc, remoteServer)
			}

			_, err = pc.WriteTo(payload, tgtUDPAddr) // accept only UDPAddr despite the signature
			if err != nil {
				logging.Error("UDP remote write error: %v", err)
				continue
			}
		}
	}()
	return nil
}

// ConnectionStage event
func (s *ShadowsocksProxy) ConnectionStage(proxy, client, target net.Addr) {
	s.MessageRoute <- record.ConnectionProxyRequest{
		ConnectionPair: record.ConnectionPair{
			ClientAddr: client,
			ProxyAddr:  proxy,
			TargetAddr: target,
		},
	}
}

func (s *ShadowsocksProxy) String() string {
	result, err := json.Marshal(s)
	if err != nil {
		log.Err(err)
		return ""
	}
	return string(result)
}

// Packet NAT table
type natmap struct {
	sync.RWMutex
	m       map[string]net.PacketConn
	timeout time.Duration
}

func newNATmap(timeout time.Duration) *natmap {
	m := &natmap{}
	m.m = make(map[string]net.PacketConn)
	m.timeout = timeout
	return m
}

func (m *natmap) Get(key string) net.PacketConn {
	m.RLock()
	defer m.RUnlock()
	return m.m[key]
}

func (m *natmap) Set(key string, pc net.PacketConn) {
	m.Lock()
	defer m.Unlock()

	m.m[key] = pc
}

func (m *natmap) Del(key string) net.PacketConn {
	m.Lock()
	defer m.Unlock()

	pc, ok := m.m[key]
	if ok {
		delete(m.m, key)
		return pc
	}
	return nil
}

func (m *natmap) Add(peer net.Addr, dst, src net.PacketConn, role mode) {
	m.Set(peer.String(), src)

	go func() {
		timedCopy(dst, peer, src, m.timeout, role)
		if pc := m.Del(peer.String()); pc != nil {
			pc.Close()
		}
	}()
}

// copy from src to dst at target with read timeout
func timedCopy(dst net.PacketConn, target net.Addr, src net.PacketConn, timeout time.Duration, role mode) error {
	buf := pool.GetUdpBuf()
	defer pool.PutUdpBuf(buf)

	for {
		src.SetReadDeadline(time.Now().Add(timeout))
		n, raddr, err := src.ReadFrom(buf)
		if err != nil {
			return errors.Cause(err)
		}

		switch role {
		case remoteServer: // server -> client: add original packet source
			srcAddr := socks.ParseAddr(raddr.String())
			copy(buf[len(srcAddr):], buf[:n])
			copy(buf, srcAddr)
			_, err = dst.WriteTo(buf[:len(srcAddr)+n], target)
		case relayClient: // client -> user: strip original packet source
			srcAddr := socks.SplitAddr(buf[:n])
			_, err = dst.WriteTo(buf[len(srcAddr):n], target)
		case socksClient: // client -> socks5 program: just set RSV and FRAG = 0
			_, err = dst.WriteTo(append([]byte{0, 0, 0}, buf[:n]...), target)
		}

		if err != nil {
			return errors.Cause(err)
		}
	}
}
