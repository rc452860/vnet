package network

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"runtime/debug"
	"strings"
	"time"
)

func NewListener(addr string,timeout time.Duration) *Listener{
	listener := new(Listener)
	listener.Timeout = timeout
	listener.addr = addr
	return listener
}

type Listener struct {
	addr    string
	Timeout time.Duration
	TCP     *net.TCPListener
}

func (l *Listener) ListenTCP(fn func(request *Request)) error {
	if l.addr == ""{
		return errors.New("listener addr is empty")
	}

	listen, err := net.Listen("tcp", l.addr)
	if err != nil {
		return err
	}
	l.TCP = listen.(*net.TCPListener)
	go func() {
		defer func() {
			if e := recover(); e != nil {
				logrus.Errorf("ListenTCP crashed , err : %s , \ntrace:%s", e, string(debug.Stack()))
			}
		}()
		for {
			err := l.TCP.SetDeadline(time.Now().Add(l.Timeout))
			if err != nil {
				logrus.Error("[%s] listen set timeout error:", err.Error())
				_ = l.Close()
				return
			}
			con, err := l.TCP.Accept()
			// TODO: https://liudanking.com/network/go-%E4%B8%AD%E5%A6%82%E4%BD%95%E5%87%86%E7%A1%AE%E5%9C%B0%E5%88%A4%E6%96%AD%E5%92%8C%E8%AF%86%E5%88%AB%E5%90%84%E7%A7%8D%E7%BD%91%E7%BB%9C%E9%94%99%E8%AF%AF/
			if err != nil {
				errString := err.Error()
				switch {
				case strings.Contains(errString, "timeout"):
					logrus.Warningf("[%s] listener accept timeout")
					continue
				default:
					fmt.Printf("[%s] listener Unknown error:%s", errString)
					return
				}
			}
			go func(){
				defer func() {
					if e := recover(); e != nil {
						logrus.WithFields(logrus.Fields{}).Errorf("connection handle crashed , err : %s , \ntrace:%s", e, string(debug.Stack()))
					}
				}()
				fn(NewRequest(true,con))
			}()
		}
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
