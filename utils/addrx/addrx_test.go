package addrx

import (
	"net"
	"testing"
)

func Benchmark_GetIPFromAddr(t *testing.B) {
	client, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:808")
	for i := 0; i < t.N; i++ {
		GetIPFromAddr(client)
	}
	t.Log(GetIPFromAddr(client))
	t.ReportAllocs()
}

func Benchmark_GetPortFromAddr(t *testing.B) {
	client, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:808")
	for i := 0; i < t.N; i++ {
		GetPortFromAddr(client)
	}
	t.Log(GetPortFromAddr(client))
	t.ReportAllocs()
}

func Benchmark_GetNetworkFromAddr(t *testing.B) {
	client, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:808")
	for i := 0; i < t.N; i++ {
		GetNetworkFromAddr(client)
	}
	t.Log(GetNetworkFromAddr(client))
	t.ReportAllocs()
}
