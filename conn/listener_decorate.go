package conn

import (
	"net"
	"time"

	"github.com/rc452860/vnet/pool"
)

type IListener interface {
	Accept() (IConn, error)
	Close() error
}

type DefaultListener struct {
	net.Listener
	closed chan struct{}
}
type Handle func(IConn)

func ListenTcp(addr string, handle Handle) (IListener, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		logging.Error(err.Error())
		return nil, err
	}
	server, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		logging.Error(err.Error())
		return nil, err
	}

	listener := DefaultListener{
		Listener: server,
		closed:   make(chan struct{}),
	}

	go func() {
		for {
			select {
			case <-listener.closed:
				return
			default:
			}
			listener.Listener.(*net.TCPListener).SetDeadline(time.Now().Add(1e9))
			con, err := listener.Accept()
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}

			if err != nil {
				logging.Error(err.Error())
				continue
			}

			go handle(con)
		}
	}()

	return listener, nil
}

func (d DefaultListener) Accept() (IConn, error) {
	con, err := d.Listener.Accept()
	if err != nil {
		return nil, err
	}
	err = con.(*net.TCPConn).SetKeepAlivePeriod(30 * time.Second)
	if err != nil {
		logging.Error(err.Error())
	}
	c, err := NewDefaultConn(con, TCP)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (d DefaultListener) Close() error {
	d.closed <- struct{}{}
	return d.Listener.Close()
}

// UDP Listener
type IUDPListener interface {
	net.PacketConn
}

type DefaultUdpListener struct {
	net.PacketConn
	closed chan struct{}
}

type UdpPacket struct {
	Buf  []byte
	Addr *net.UDPAddr
	net.PacketConn
}

func ListenUdp(addr string, handle func(*UdpPacket)) (IUDPListener, error) {
	l, err := net.ListenPacket("udp", addr)
	if err != nil {
		return nil, err
	}
	closed := make(chan struct{})
	listener := &DefaultUdpListener{
		PacketConn: l,
		closed:     closed,
	}
	for {
		buf := pool.GetUdpBuf()
		n, raddr, err := l.ReadFrom(buf)
		if err != nil {
			logging.Err(err)
			continue
		}
		logging.Debug("recive packet: %v", raddr)

		udpPacket := &UdpPacket{
			Buf:        buf[:n],
			Addr:       raddr.(*net.UDPAddr),
			PacketConn: l,
		}
		go func() {
			defer pool.PutUdpBuf(buf)
			handle(udpPacket)
		}()
	}
	return listener, nil
}

// func ListenUdpS(password, method, addr string, handle func(*UdpPacket)) (IUDPListener, error) {
// 	l, err := net.ListenPacket("udp", addr)
// 	if err != nil {
// 		return nil, err
// 	}
// 	l, err = ciphers.CipherPacketDecorate(password, method, l)
// 	if err != nil {
// 		logging.Error("can't init packet decorate")
// 		return nil, err
// 	}
// 	closed := make(chan struct{})
// 	listener := &DefaultUdpListener{
// 		PacketConn: l,
// 		closed:     closed,
// 	}
// 	for {
// 		buf := pool.GetUdpBuf()
// 		n, raddr, err := l.ReadFrom(buf)
// 		if err != nil {
// 			logging.Err(err)
// 			continue
// 		}
// 		logging.Debug("recive packet: %v", raddr)

// 		udpPacket := &UdpPacket{
// 			Buf:        buf[:n],
// 			Addr:       raddr.(*net.UDPAddr),
// 			PacketConn: l,
// 		}
// 		go func() {
// 			defer pool.PutUdpBuf(buf)
// 			handle(udpPacket)
// 		}()
// 	}
// 	return listener, nil
// }
