package server

import (
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/rc452860/vnet/utils/datasize"

	"github.com/rc452860/vnet/ciphers"
	"github.com/rc452860/vnet/conn"
	"github.com/rc452860/vnet/log"
	"github.com/rc452860/vnet/pool"
	"github.com/rc452860/vnet/proxy"
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
	Method              string        `json:"method,omitempty"`
	Password            string        `json:"password,omitempty"`
	Port                int           `json:"port,omitempty"`
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

func (s *ShadowsocksProxy) ConfigLimit(limit uint64) {
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
		logging.Info("server port: %v limit is: %v", s.Port, trafficLimit)
		s.ReadLimit = rate.NewLimiter(rate.Limit(trafficLimit), int(trafficLimit))
		s.WriteLimit = rate.NewLimiter(rate.Limit(trafficLimit), int(trafficLimit))
	}
	return nil

}

func (s *ShadowsocksProxy) ConfigTimeout(timeout time.Duration) error {
	if timeout == 0 {
		s.TCPTimeout = 3e9
		s.UDPTimeout = 3e9
	} else {
		s.TCPTimeout = timeout
		s.UDPTimeout = timeout
	}
	log.Info("%s:%v timeout:%v", s.Host, s.Port, s.TCPTimeout)
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
	go s.ProxyService.TrafficMeasure()
	return nil
}

// Stop proxy
func (s *ShadowsocksProxy) Stop() error {
	return s.ProxyService.Stop()
}

// statistics upload traffic
func (s *ShadowsocksProxy) upload(con conn.IConn, up int) {
	s.ProxyService.TrafficMQ <- proxy.TrafficMessage{
		Network: "tcp",
		LAddr:   con.LocalAddr().String(),
		RAddr:   con.RemoteAddr().String(),
		UpBytes: uint64(up),
	}
}

// statics download traffic
func (s *ShadowsocksProxy) download(con conn.IConn, down int) {
	s.ProxyService.TrafficMQ <- proxy.TrafficMessage{
		Network:   "tcp",
		LAddr:     con.LocalAddr().String(),
		RAddr:     con.RemoteAddr().String(),
		DownBytes: uint64(down),
	}
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
	logging.Info("listening TCP on %s", addr)
	if err != nil {
		logging.Error(err.Error())
		return err
	}
	s.Tcp = server

	go func() {
		defer server.Close()
		for {
			// select {
			// case <-s.ProxyService.TcpClose:
			// 	s.Tcp.Close()
			// 	s.Tcp = nil
			// 	return
			// default:
			// }
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
				/** 流量记录装饰器 */
				lcd, err = conn.TrafficDecorate(lcd, s.upload, s.download)
				if err != nil {
					logging.Err(err)
					return
				}
				/** 限流装饰器 */
				lcd, _ = conn.TrafficLimitDecorate(lcd, s.ReadLimit, s.WriteLimit)

				// lcd, err = conn.TimerDecorate(lcd, s.TcpTimeout, s.TcpTimeout)
				// if err != nil {
				// 	logging.Err(err)
				// 	return
				// }
				/** 加密装饰器 */
				lcd, err = ciphers.CipherDecorate(s.Password, s.Method, lcd)
				if err != nil {
					logging.Err(err)
					return
				}

				/** 读取目标地址 */
				targetAddr, err := socks.ReadAddr(lcd)
				if err != nil {
					logging.Err(err)
					return
				}

				rc, err := net.Dial("tcp", targetAddr.String())
				if err != nil {
					logging.Error("connect target:%s error", targetAddr)
					logging.Err(err)
					return
				}
				defer rc.Close()

				rc.(*net.TCPConn).SetKeepAlive(true)
				logging.Info("tcp %s <----> %s", lcd.RemoteAddr(), targetAddr)

				/** 默认装饰器 */
				rcd, err := conn.DefaultDecorate(rc, conn.TCP)
				if err != nil {
					logging.Err(err)
					return
				}

				// rcd, _ = conn.TrafficLimitDecorate(rcd, s.ReadLimit, s.WriteLimit)

				/** 流量统计装饰器 */
				// rcd, err = conn.TrafficDecorate(rcd)
				// if err != nil {
				// 	logging.Err(err)
				// 	return
				// }

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
	return n, rs.N, err
}

// udp upload traffic count
func (s *ShadowsocksProxy) udpUpload(laddr, raddr string, n int) {
	s.ProxyService.TrafficMQ <- proxy.TrafficMessage{
		Network: "tcp",
		LAddr:   laddr,
		RAddr:   raddr,
		UpBytes: uint64(n),
	}
}

// udp download traffic count
func (s *ShadowsocksProxy) udpDownload(laddr, raddr string, n int) {
	s.ProxyService.TrafficMQ <- proxy.TrafficMessage{
		Network:   "tcp",
		LAddr:     laddr,
		RAddr:     raddr,
		DownBytes: uint64(n),
	}
}

// Listen on addr for encrypted packets and basically do UDP NAT.
func (s *ShadowsocksProxy) startUDP() error {
	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)
	server, err := net.ListenPacket("udp", addr)
	if err != nil {
		logging.Error("UDP remote listen error: %v", err)
		return err
	}
	s.Udp = server
	server, err = ciphers.CipherPacketDecorate(s.Password, s.Method, server)
	server = conn.PacketTrafficConnDecorate(server, s.udpUpload, s.udpDownload)
	if err != nil {
		logging.Error("UDP CipherPacketDecorate init error: %v", err)
		return err
	}

	nm := newNATmap(s.UDPTimeout)
	buf := pool.GetUdpBuf()
	defer pool.PutUdpBuf(buf)

	logging.Info("listening UDP on %s", addr)

	go func() {
		defer server.Close()
		for {
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
				logging.Error("failed to split target address from packet: %v", buf[:n])
				continue
			}
			logging.Info("udp %s <----> %s", raddr, tgtAddr)

			tgtUDPAddr, err := net.ResolveUDPAddr("udp", tgtAddr.String())
			if err != nil {
				logging.Error("failed to resolve target UDP address: %v", err)
				continue
			}

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
			return err
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
			return err
		}
	}
}
