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
}

func Benchmark_OriginFormat(t *testing.B) {
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		time.Now().Format("06")
	}
}
