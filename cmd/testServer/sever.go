package main

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/rc452860/vnet/common/config"
	"github.com/rc452860/vnet/common/log"
	"github.com/rc452860/vnet/common/pool"
	"github.com/rc452860/vnet/network/conn"
	"github.com/rc452860/vnet/proxy/client"
	"github.com/rc452860/vnet/proxy/server"
	"github.com/rc452860/vnet/service"
	thttp "github.com/rc452860/vnet/testing/servers/http"
)

const number = 10

func main() {
	config.LoadConfig("config.json")
	// start pprof
	go func() {
		http.ListenAndServe("0.0.0.0:6060", nil)
	}()

	for i := 10000; i < 10000+number; i++ {
		service.CurrentShadowsocksService().Add("0.0.0.0", "aes-128-gcm", "killer", i, server.ShadowsocksArgs{
			ConnectTimeout: 3000,
			Limit:          0,
			TCPSwitch:      "",
			UDPSwitch:      "",
		})
		err := service.CurrentShadowsocksService().Start(i)
		if err != nil {
			log.Err(err)
		}
	}
	log.Info("all service is started")

	Testing()
	time.Sleep(10 * time.Second)
	log.Info("=====================================================================")
	ch := make(chan os.Signal, 2)

	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch
	gw := new(sync.WaitGroup)
	for i := 10000; i < 10000+number; i++ {
		gw.Add(1)
		go func(index int) {
			gw.Done()
			service.CurrentShadowsocksService().Stop(index)
		}(i)
	}
	gw.Wait()
	time.Sleep(3 * time.Second)
	log.Info("all service stoped")
	<-ch
	log.Info("bybe~")
}

func Testing() {
	log.GetLogger("root").Level = log.INFO
	thttp.StartFakeFileServer()
	wg := new(sync.WaitGroup)
	for i := 10000; i < 20000; i++ {
		wg.Add(1)
		go func() {
			TestingClient(int(10000 + 1))
			wg.Done()
		}()
	}
	wg.Wait()
	// runtime.GC()
}

func TestingClient(i int) {
	s, c := net.Pipe()
	defer s.Close()
	defer c.Close()
	client := client.NewShadowsocksClient("127.0.0.1", "aes-128-gcm", "killer", i)
	cs, err := conn.NewDefaultConn(c, "pipe")
	go client.TcpProxy(cs, "127.0.0.1", 8080)
	var httpRequest bytes.Buffer
	httpRequest.WriteString("GET /download?size=4KB HTTP/1.1\n")
	httpRequest.WriteString("Host: localhost:8080\n")
	httpRequest.WriteString("Connection: keep-alive\n")
	// httpRequest.WriteString("Connection: close\n")
	httpRequest.WriteString("Pragma: no-cache\n")
	httpRequest.WriteString("Cache-Control: no-cache\n")
	httpRequest.WriteString("Upgrade-Insecure-Requests: 1\n")
	httpRequest.WriteString("User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.80 Safari/537.36\n")
	httpRequest.WriteString("DNT: 1\n\n")
	data := httpRequest.Bytes()
	s.SetWriteDeadline(time.Now().Add(3 * time.Second))
	_, err = s.Write(data)
	if err != nil {
		return
	}
	var count int64 = 0
	buf := pool.GetBuf()
	defer pool.PutBuf(buf)
	request, err := http.ReadRequest(bufio.NewReader(bufio.NewReader(&httpRequest)))
	response, err := http.ReadResponse(bufio.NewReader(s), request)
	log.Info("content lenght:%v", response.ContentLength)
	if err != nil {
		log.Err(err)
		return
	}
	for {
		n, err := response.Body.Read(buf)
		// buff.Write(buf)
		count = count + int64(n)
		if err != nil && err != io.EOF {
			log.Err(err)
			break
		}
		if count == response.ContentLength {
			break
		}
	}
	log.Info("done %v", i)
}
