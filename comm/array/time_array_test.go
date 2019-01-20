package array

import (
	"testing"
	"time"
)

func Benchmark_Time_Array(t *testing.B) {
	timeArray := NewTimeArray(time.Second*2, true)
	for i := 0; i < t.N; i++ {
		timeArray.Add("123")
	}
}

func Benchmark_Time_Array_Remove(t *testing.B) {
	timeArray := NewTimeArray(time.Second*2, true)
	for i := 0; i < t.N; i++ {
		timeArray.Add("123")
	}
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		timeArray.Remove("123")
	}
	t.ReportAllocs()
}

func Test_Time_Array(t *testing.T) {
	timeArray := NewTimeArray(time.Second*2, true)
	tick := time.Tick(time.Second)
	index := 0
	for {
		<-tick
		index++
		if index > 10 {
			break
		}

		timeArray.Add("aaa")
	}
	timeArray.Range(func(i int, v interface{}) {
		t.Log(v.(string))
	})
}
