package utils

import (
	"os"
	"testing"
	"time"
)

func Test_IsFileExist(t *testing.T) {
	if IsFileExist("./io.go") {
		t.Log("exist")
	} else {
		t.Error("not exist")
	}

	if IsFileExist("./io_.go") {
		t.Error("exist")
	} else {
		t.Log("not exist")
	}
}

func Test_OpenFileWrite(t *testing.T) {
	file, err := os.OpenFile("aaa.txt", os.O_APPEND|os.O_CREATE, 0666)
	defer file.Close()
	if err != nil {
		panic(err)
	}
	for i := 0; i < 20; i++ {
		file.WriteString("aaa\n")
		time.Sleep(1 * time.Second)
	}
}

func Benchmark_OpenFile(t *testing.B) {
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		file, _ := os.OpenFile("aaa.txt", os.O_APPEND|os.O_CREATE, 0666)
		file.Close()
	}

}

func Benchmark_IsFileExist(t *testing.B) {
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		IsFileExist("./io.go")
	}
}
