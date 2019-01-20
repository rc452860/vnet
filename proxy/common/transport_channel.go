package common

import (
	"io"
	"strings"

	"github.com/rc452860/vnet/pool"

	"github.com/rc452860/vnet/conn"
	"github.com/rc452860/vnet/comm/log"
)

type TcpChannel struct{}

func Recover(fs ...func()) {
	if err := recover(); err != nil {
		log.Error("[PANIC] %v", err)
		for _, f := range fs {
			f()
		}
	}
}

func (d *TcpChannel) Transport(lc, sc conn.IConn) {
	errChan := make(chan error, 2)
	go func() {
		defer Recover(func() {
			select {
			case errChan <- nil:
			default:
			}
		})
		d.send(sc, lc, errChan)
		<-errChan
	}()
	go func() {
		defer Recover(func() {
			select {
			case errChan <- nil:
			default:
			}
		})
		d.send(lc, sc, errChan)

	}()
	<-errChan
	lc.Close()
	sc.Close()
}

func (d *TcpChannel) send(from, to conn.IConn, errChan chan error) {
	var (
		buf []byte
		n   int
		err error
	)
	buf = pool.GetBuf()
	defer pool.PutBuf(buf)
	for {
		n, err = from.Read(buf)
		// @fix 空数据返回引发断连
		//if n == 0 {
		//errChan <- nil
		//return
		//}
		if err != nil {
			if err != io.EOF && !strings.Contains(err.Error(), "use of closed network connection") {
				log.Error("[ID:%d] [DirectChannel] DirectChannel Transport: %v", from.GetID(), err)
			}
			errChan <- err
			return
		}
		n, err = to.Write(buf[:n])
		if err != nil {
			if err != io.EOF && !strings.Contains(err.Error(), "use of closed network connection") {
				log.Error("[ID:%d] [DirectChannel] DirectChannel Transport: %v", to.GetID(), err)
			}
			errChan <- err
			return
		}
	}
}
