package server

import (
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/rc452860/vnet/common/pool"

	"github.com/rc452860/vnet/network/ciphers"
	"github.com/rc452860/vnet/network/conn"
	"github.com/rc452860/vnet/socks"
)

// mockUDPServer is mock udp server for shadowsocks udp test
func mockUDPServer(t *testing.T) {
	conn, err := net.ListenPacket("udp", "0.0.0.0:8081")
	if err != nil {
		t.Error(err)
		t.FailNow()
		return
	}
	buf := pool.GetUdpBuf()
	n, addr, err := conn.ReadFrom(buf)
	if "hello" == string(buf[:n]) {
		t.Logf(string(buf[:n]))
		t.Logf("udp success")
	}
	n, _ = conn.WriteTo([]byte("hello"), addr)
	conn.Close()
}
func Test_NewServer(t *testing.T) {
	testShadowsocksProxy(t, "0.0.0.0", "rc4-md5", "killer", 8080)
	t.Logf("--------------------rc4-md5 success--------------------")
	testShadowsocksProxy(t, "0.0.0.0", "aes-128-cfb", "killer", 8080)
	t.Logf("--------------------aes-128-cfb success--------------------")

}

func testShadowsocksProxy(t *testing.T, host, method, password string, port int) {
	ss, _ := NewShadowsocks(host, method, password, port, ShadowsocksArgs{
		Limit:          4096 * 1024,
		ConnectTimeout: 3 * time.Second,
	})
	go ss.Start()
	time.Sleep(1 * time.Second)
	transport := &http.Transport{
		Proxy: nil,
		Dial: func(network, addr string) (net.Conn, error) {
			con, _ := net.Dial("tcp", fmt.Sprintf("%s:%v", host, port))
			c, _ := conn.DefaultDecorate(con, conn.TCP)
			c, err := ciphers.CipherDecorate(password, method, c)
			c.Write(socks.ParseAddr(addr).Raw)
			return c, err
		},
		TLSHandshakeTimeout: 10 * time.Second,
	}
	client := &http.Client{
		Transport: transport,
	}
	response, err := client.Get("http://baidu.com")
	if err != nil {
		t.Error(err)
		return
	}
	if response.StatusCode < 200 && response.StatusCode > 400 {
		t.Fatal("http status error")
	}

	t.Logf("tcp success: %s", ss.String())
	go mockUDPServer(t)
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%v", "127.0.0.1", port))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	packet, err := net.ListenPacket("udp", "0.0.0.0:12345")
	packet, err = ciphers.CipherPacketDecorate(password, method, packet)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	sendBuf := pool.GetBuf()
	targetAddr := socks.ParseAddr("127.0.0.1:8081")
	copy(sendBuf, targetAddr.Raw)
	copy(sendBuf[len(targetAddr.Raw):], []byte("hello"))
	n, err := packet.WriteTo(sendBuf[:len(targetAddr.Raw)+5], addr)
	pool.PutBuf(sendBuf)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	buf := pool.GetBuf()
	n, _, err = packet.ReadFrom(buf)
	t.Logf("%v", n)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fromAddr := socks.SplitAddr(buf[:n])
	if string(buf[len(fromAddr.Raw):n]) != "hello" {
		t.Error("recive is not compare hello")
	}
	packet.Close()
	ss.Stop()
}
