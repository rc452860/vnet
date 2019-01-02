package log

import (
	"fmt"
	"testing"
)

func Test_GetRuntimeInfo(t *testing.T) {
	funcName, file, line := GetRuntimeInfoShortFormat()
	fmt.Print("-----")
	fmt.Print(funcName, "\n", file, "\n", line)
}

func Benchmark_GetRuntimeInfoShortFormat(t *testing.B) {
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		GetRuntimeInfoShortFormat()
	}
}

func Benchmark_GetRuntimeInfo(t *testing.B) {
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		GetRuntimeInfo()
	}
}
