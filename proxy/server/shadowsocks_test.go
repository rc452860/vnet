package server

import (
	"io"
	"net"
	"testing"
	"time"

	"github.com/rc452860/vnet/ciphers"
	"github.com/rc452860/vnet/conn"
	"github.com/rc452860/vnet/pool"
	"github.com/rc452860/vnet/socks"
)

func Test_NewServer(t *testing.T) {
	go NewServer("0.0.0.0", "aes-128-gcm", "killer", 8080)
	time.Sleep(3 * time.Second)
	con, _ := net.Dial("tcp", "0.0.0.0:8080")
	c, _ := conn.DefaultDecorate(con, conn.TCP)
	c, err := ciphers.CipherDecorate("killer", "aes-128-gcm", c)
	if err != nil {
		logging.Error(err.Error())
	}
	c.Write(socks.ParseAddr("baidu.com.com:80"))
	c.Write([]byte("GET / HTTP/1.1\n"))
	c.Write([]byte("host: baidu.com\n"))
	c.Write([]byte("Connection: keep-alive\n"))
	c.Write([]byte("User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.80 Safari/537.36\n"))
	c.Write([]byte("Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8\n"))
	c.Write([]byte("Accept-Encoding: gzip, deflate, br\n"))
	c.Write([]byte("Accept-Language: zh-CN,zh;q=0.9\n"))

	c.Write([]byte("\n\n"))
	results := pool.GetBuf()

	for {
		n, err := c.Read(results)
		if err != nil && err == io.EOF {
			logging.Info("read end")
			return
		}
		if err != nil {
			logging.Error(err.Error())
		}
		logging.Info(string(results[0:n]))
		if n <= 0 {
			return
		}
	}
}
