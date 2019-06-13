package network

import (
	"errors"
	"github.com/sirupsen/logrus"
	"net"
	"runtime/debug"
	"strings"
	"time"
)

func NewListener(addr string, timeout time.Duration) *Listener {
	listener := new(Listener)
	listener.Timeout = timeout
	listener.Addr = addr
	return listener
}

type Listener struct {
	Addr    string
	Timeout time.Duration
	TCP     *net.TCPListener
	UDP     net.PacketConn
}

func (l *Listener) ListenTCP(fn func(request *Request)) error {
	if l.Addr == "" {
		return errors.New("listener Addr is empty")
	}

	listen, err := net.Listen("tcp", l.Addr)
	if err != nil {
		return err
	}
	logrus.Infof("Listener listen on: %s", l.Addr)
	l.TCP = listen.(*net.TCPListener)
	go func() {
		defer func() {
			if e := recover(); e != nil {
				logrus.Errorf("ListenTCP crashed , err : %s , \ntrace:%s", e, string(debug.Stack()))
			}
		}()
		for {
			if l.Timeout != 0 {
				err := l.TCP.SetDeadline(time.Now().Add(l.Timeout))
				if err != nil {
					logrus.Error("[%s] listen set timeout error:", err.Error())
					_ = l.Close()
					return
				}
			}
			con, err := l.TCP.Accept()
			// TODO: https://liudanking.com/network/go-%E4%B8%AD%E5%A6%82%E4%BD%95%E5%87%86%E7%A1%AE%E5%9C%B0%E5%88%A4%E6%96%AD%E5%92%8C%E8%AF%86%E5%88%AB%E5%90%84%E7%A7%8D%E7%BD%91%E7%BB%9C%E9%94%99%E8%AF%AF/
			if err != nil {
				errString := err.Error()
				switch {
				case strings.Contains(errString, "timeout"):
					continue
				default:
					logrus.Errorf("[%s] listener Unknown error:%s", errString)
					return
				}
			}
			go func() {
				defer func() {
					if e := recover(); e != nil {
						logrus.WithFields(logrus.Fields{}).Errorf("connection handle crashed , err : %s , \ntrace:%s", e, string(debug.Stack()))
					}
				}()
				fn(NewRequestWithTCP(con))
			}()
		}
	}()
	return nil
}

func (l *Listener) ListenUDP(fn func(request *Request)) error {
	if l.Addr == "" {
		return errors.New("listener Addr is empty")
	}

	listen, err := net.ListenPacket("udp", l.Addr)
	if err != nil{
		return err
	}
	l.UDP = listen
	go func(){
		defer func() {
			if e := recover(); e != nil {
				logrus.Errorf("ListenUDP crashed , err : %s , \ntrace:%s", e, string(debug.Stack()))
			}
		}()
		fn(NewRequestWithUDP(l.UDP))
	}()
	return nil

}

func (l *Listener) Close() error {
	err := l.TCP.Close()
	if err != nil {
		return err
	}
	return nil
}
