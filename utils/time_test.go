package utils

import (
	"fmt"
	"testing"
	"time"
)

func Test_Format(t *testing.T) {
	fmt.Println(Format("yyyy-MM-dd HH:mm:ss", time.Now()))
}

func Benchmark_Format(t *testing.B) {
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		Format("yyyy-MM-dd HH:mm:ss", time.Now())
	}
	t.Log(Format("yyyy-MM-dd HH:mm:ss", time.Now()))
	t.ReportAllocs()
}

func Benchmark_format_string(t *testing.B) {
	a := []byte{0x45, 0x46}
	t.ResetTimer()
	for i := 1; i < t.N; i++ {
		_ = string(a)
	}
	t.ReportAllocs()
}

func Benchmark_OriginFormat(t *testing.B) {
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		time.Now().Format("2006-01-02 15:04:05")
	}
	t.ReportAllocs()
}
