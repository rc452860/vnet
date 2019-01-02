package server

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/rc452860/vnet/config"
	"github.com/rc452860/vnet/pool"
	"github.com/rc452860/vnet/proxy"

	"github.com/rc452860/vnet/socks"

	"github.com/rc452860/vnet/ciphers"

	"github.com/rc452860/vnet/conn"

	"github.com/rc452860/vnet/log"
)

type mode int

const (
	remoteServer mode = iota
	relayClient
	socksClient
)

var logging *log.Logging

func init() {
	logging = log.GetLogger("root")
}

type ShadowsocksProxy struct {
	proxy.ProxyService
	Host       string
	Method     string
	Password   string
	Port       int
	TcpTimeout time.Duration
	UdpTimeout time.Duration
}

func NewShadowsocsk(host string, method string, password string, port int) *ShadowsocksProxy {
	ss := &ShadowsocksProxy{
		ProxyService: proxy.NewProxyService(),
		Host:         host,
		Method:       method,
		Password:     password,
		Port:         port,
	}
	if config.CurrentConfig().ShadowsocksOptions.TcpTimeout == 0 {
		ss.TcpTimeout = 3e9
	} else {
		ss.TcpTimeout = config.CurrentConfig().ShadowsocksOptions.TcpTimeout
	}

	if config.CurrentConfig().ShadowsocksOptions.UdpTimeout == 0 {
		ss.UdpTimeout = 3e9
	} else {
		ss.UdpTimeout = config.CurrentConfig().ShadowsocksOptions.UdpTimeout
	}
	return ss
}

func (this ShadowsocksProxy) Start() error {
	if err := this.startTcp(); err != nil {
		return err
	}
	if err := this.startUdp(); err != nil {
		return err
	}
	return nil
}

func (this ShadowsocksProxy) Stop() error {
	return this.ProxyService.Stop()
}

// start shadowsocks tcp proxy service
func (this ShadowsocksProxy) startTcp() error {
	addr := fmt.Sprintf("%s:%d", this.Host, this.Port)
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
	this.Tcp = server

	go func() {
		defer server.Close()
		for {
			select {
			case <-this.TcpClose:
				this.Tcp = nil
				return
			default:
			}
			server.SetDeadline(time.Now().Add(this.TcpTimeout))
			lcon, err := server.Accept()
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}

			if err != nil {
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
				lcd, err = conn.TrafficDecorate(lcd)
				if err != nil {
					logging.Err(err)
					return
				}

				// /** 超时装饰器 */
				// lcd, err = conn.TimerDecorate(lcd, this.TcpTimeout, this.TcpTimeout)
				// if err != nil {
				// 	logging.Err(err)
				// 	return
				// }
				/** 加密装饰器 */
				lcd, err = ciphers.CipherDecorate(this.Password, this.Method, lcd)
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
				logging.Info("%s <----> %s", lcd.RemoteAddr(), targetAddr)

				/** 默认装饰器 */
				rcd, err := conn.DefaultDecorate(rc, conn.TCP)
				if err != nil {
					logging.Err(err)
					return
				}

				/** 流量统计装饰器 */
				// rcd, err = conn.TrafficDecorate(rcd)
				// if err != nil {
				// 	logging.Err(err)
				// 	return
				// }

				_, _, err = relayTcp(lcd, rcd)
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
func relayTcp(left, right net.Conn) (int64, int64, error) {
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

// Listen on addr for encrypted packets and basically do UDP NAT.
func (this ShadowsocksProxy) startUdp() error {
	addr := fmt.Sprintf("%s:%d", this.Host, this.Port)
	server, err := net.ListenPacket("udp", addr)
	if err != nil {
		logging.Error("UDP remote listen error: %v", err)
		return err
	}
	this.Udp = server
	server, err = ciphers.CipherPacketDecorate(this.Password, this.Method, server)
	if err != nil {
		logging.Error("UDP CipherPacketDecorate init error: %v", err)
		return err
	}

	nm := newNATmap(this.UdpTimeout)
	buf := pool.GetUdpBuf()
	defer pool.PutUdpBuf(buf)

	logging.Info("listening UDP on %s", addr)

	go func() {
		defer server.Close()
		for {
			select {
			case <-this.UdpClose:
				this.Udp = nil
				return
			default:
			}
			server.SetDeadline(time.Now().Add(this.UdpTimeout))
			n, raddr, err := server.ReadFrom(buf)
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			if err != nil {
				logging.Error("UDP remote read error: %v", err)
				continue
			}

			tgtAddr := socks.SplitAddr(buf[:n])
			if tgtAddr == nil {
				logging.Error("failed to split target address from packet: %q", buf[:n])
				continue
			}

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
