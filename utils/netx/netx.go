package netx

import (
	"github.com/pkg/errors"
	"github.com/rc452860/vnet/common/log"
	"github.com/rc452860/vnet/common/pool"
	"github.com/rc452860/vnet/utils/goroutine"
	"github.com/rc452860/vnet/utils/socksproxy"
	"io"
	"net"
	"sync"
	"time"
)

// DuplexCopyTcp will return 3 result
// up means left connection to right connection transfer data count
// down means right connection to left connections transfer data count
// and the last result is error
func DuplexCopyTcp(left, right net.Conn) (up, down int64, err error) {
	type res struct {
		N   int64
		Err error
	}
	ch := make(chan res)
	defer func() {
		if e := recover(); e != nil {
			log.Error("panic in timedCopy: %v", e)
		}
	}()

	go goroutine.Protect(func() {
		n, err := io.Copy(right, left)
		_ = right.SetDeadline(time.Now()) // wake up the other goroutine blocking on right
		_ = left.SetDeadline(time.Now())  // wake up the other goroutine blocking on left
		ch <- res{n, err}
	})

	up, err = io.Copy(left, right)
	_ = right.SetDeadline(time.Now()) // wake up the other goroutine blocking on right
	_ = left.SetDeadline(time.Now())  // wake up the other goroutine blocking on left
	rs := <-ch

	if err == nil {
		err = rs.Err
	}
	return up, rs.N, errors.Cause(err)
}

// Packet NAT table
type NatMap struct {
	sync.RWMutex
	m       map[string]net.PacketConn
	timeout time.Duration
}

func NewNatMap(timeout time.Duration) *NatMap {
	m := &NatMap{}
	m.m = make(map[string]net.PacketConn)
	m.timeout = timeout
	return m
}

func (m *NatMap) Get(key string) net.PacketConn {
	m.RLock()
	defer m.RUnlock()
	return m.m[key]
}

func (m *NatMap) Set(key string, pc net.PacketConn) {
	m.Lock()
	defer m.Unlock()
	m.m[key] = pc
}

func (m *NatMap) Del(key string) net.PacketConn {
	m.Lock()
	defer m.Unlock()

	pc, ok := m.m[key]
	if ok {
		delete(m.m, key)
		return pc
	}
	return nil
}

func (m *NatMap) Add(peer net.Addr, dst, src net.PacketConn) {
	m.Set(peer.String(), src)
	go goroutine.Protect(func() {
		_ = timedCopy(dst, peer, src, m.timeout)
		if pc := m.Del(peer.String()); pc != nil {
			_ = pc.Close()
		}
	})
}

// copy from src to dst at target with read timeout
func timedCopy(dst net.PacketConn, target net.Addr, src net.PacketConn, timeout time.Duration) error {
	buf := pool.GetBuf()
	defer pool.PutBuf(buf)
	defer func() {
		if e := recover(); e != nil {
			log.Error("panic in timedCopy: %v", e)
		}
	}()

	for {
		_ = src.SetReadDeadline(time.Now().Add(timeout))
		n, raddr, err := src.ReadFrom(buf)
		if err != nil {
			return errors.Cause(err)
		}

		srcAddr := socksproxy.ParseAddr(raddr.String())
		srcAddrByte := srcAddr.Raw
		copy(buf[len(srcAddrByte):], buf[:n])
		copy(buf, srcAddrByte)
		_, err = dst.WriteTo(buf[:len(srcAddrByte)+n], target)

		if err != nil {
			return errors.Cause(err)
		}
	}
}

