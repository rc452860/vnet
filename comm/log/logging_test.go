package log

import (
	"strconv"
	"strings"
	"testing"
	"time"
)

func Test_Logging(t *testing.T) {
	logging := GetLogger("root", INFO)

	logging.Debug("aaa")
	logging.Info("bbb")
	logging.Warn("bbb")
	logging.Error("bbb")
}

func Benchmark_Logging(t *testing.B) {
	logging := GetLogger("root", INFO)

	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		logging.Info("bbb")
	}
	t.ReportAllocs()
}

func Benchmark_Format(t *testing.B) {
	pattern := PatternLogFormatterFactory(pattern)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		pattern.Format("bbbb", "INFO")
	}
}

func Benchmark_Replace(t *testing.B) {
	replacer := strings.NewReplacer(
		LEVEL,
		"INFO",
		TIME,
		time.Now().Format("2006-01-02 15:04:05"),
		FILE,
		"file",
		FUNC,
		"funcName",
		LINENO,
		strconv.Itoa(1),
		MESSAGE,
		"aaa",
	)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		replacer.Replace(pattern)
	}
	t.ReportAllocs()
}
