package main

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/rc452860/vnet/proxy/client"

	"github.com/rc452860/vnet/common/config"
	"github.com/rc452860/vnet/common/log"
	"github.com/rc452860/vnet/common/pool"
	"github.com/rc452860/vnet/network/conn"
	"github.com/rc452860/vnet/proxy/server"
	thttp "github.com/rc452860/vnet/testing/servers/http"
	"github.com/rc452860/vnet/utils/datasize"
)

func main() {
	config.LoadConfig("config.json")
	log.GetLogger("root").Level = log.INFO
	proxy, err := server.NewShadowsocks("0.0.0.0", "aes-128-gcm", "killer", 1090, server.ShadowsocksArgs{
		Limit:          4 * 1024 * 1024,
		ConnectTimeout: 0,
	})
	proxy.Start()
	if err != nil {
		log.Err(err)
		return
	}
	time.Sleep(time.Second)
	thttp.StartFakeFileServer()
	s, c := net.Pipe()
	client := client.NewShadowsocksClient("127.0.0.1", "aes-128-gcm", "killer", 1090)
	cs, err := conn.NewDefaultConn(c, "pipe")
	go client.TcpProxy(cs, "localhost", 8080)
	var httpRequest bytes.Buffer
	httpRequest.WriteString("GET /download?size=4MB HTTP/1.1\n")
	httpRequest.WriteString("Host: localhost:8080\n")
	httpRequest.WriteString("Connection: keep-alive\n")
	// httpRequest.WriteString("Connection: close\n")
	httpRequest.WriteString("Pragma: no-cache\n")
	httpRequest.WriteString("Cache-Control: no-cache\n")
	httpRequest.WriteString("Upgrade-Insecure-Requests: 1\n")
	httpRequest.WriteString("User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.80 Safari/537.36\n")
	httpRequest.WriteString("DNT: 1\n\n")
	data := httpRequest.Bytes()
	s.Write(data)
	var count int64 = 0
	buf := pool.GetBuf()
	defer pool.PutBuf(buf)
	size, _ := datasize.Parse("4MB")
	// var buff bytes.Buffer
	log.Info("%v", size)
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
	log.Info("count %v", count)
	log.Info("upload %v", proxy.UpBytes)
	// log.Info("%s", buff.Bytes()[:255])
	upspeed, _ := datasize.HumanSize(uint64(proxy.UpSpeed))
	downspeed, _ := datasize.HumanSize(uint64(proxy.DownSpeed))
	up, _ := datasize.HumanSize(proxy.UpBytes)
	down, _ := datasize.HumanSize(proxy.DownBytes)
	log.Info("upspeed:%s - downspeed:%s | up:%s - down:%s", upspeed, downspeed, up, down)
}
