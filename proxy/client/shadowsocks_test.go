package client

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rc452860/vnet/config"
	"github.com/rc452860/vnet/conn"
	"github.com/rc452860/vnet/comm/log"
	"github.com/rc452860/vnet/pool"
	"github.com/rc452860/vnet/proxy/server"

	"github.com/rc452860/vnet/utils/datasize"
)

type FakeFile struct {
	Size   int64
	Offset int64
}

func NewFakeFile(size int64) *FakeFile {
	return &FakeFile{
		Size:   size,
		Offset: 0,
	}
}
func (f *FakeFile) Read(p []byte) (n int, err error) {
	if f.Offset >= f.Size {
		return 0, io.EOF
	}
	remain := f.Size - f.Offset
	if int64(cap(p)) < remain {
		n, err = rand.Read(p)
	} else {
		n, err = rand.Read(p[:int(remain)])
	}

	f.Offset = f.Offset + int64(n)
	return n, err
}

func (f *FakeFile) Seek(offset int64, whence int) (int64, error) {
	if whence == io.SeekEnd {
		f.Offset = f.Size + offset

	}
	if whence == io.SeekStart {
		f.Offset = 0 + offset
	}
	if whence == io.SeekCurrent {
		f.Offset = f.Offset + offset
	}
	if f.Offset > f.Size || f.Offset < 0 {
		return 0, fmt.Errorf("offset is out of bounds")
	}
	return f.Offset, nil
}

func startFakeFileServer() {
	r := gin.Default()
	r.GET("download", func(c *gin.Context) {
		sizeStr := c.Query("size")
		size, err := datasize.Parse(sizeStr)
		if err != nil {
			log.Err(err)
		}
		file := NewFakeFile(int64(size))
		c.Writer.Header().Add("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.test"`, sizeStr))
		http.ServeContent(c.Writer, c.Request, fmt.Sprintf("%s.test", sizeStr), time.Now(), file)
	})
	r.Run(":8080")
}

func Test_Limit(t *testing.T) {
	config.LoadConfig("config.json")
	log.GetLogger("root").Level = log.INFO
	proxy, err := server.NewShadowsocks("0.0.0.0", "aes-128-gcm", "killer", 1090, "4MB", 0)
	proxy.Start()
	if err != nil {
		t.Error(err)
		return
	}
	time.Sleep(time.Second)
	go startFakeFileServer()
	s, c := net.Pipe()
	client := NewShadowsocksClient("127.0.0.1", "aes-128-gcm", "killer", 1090)
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
