package main

import "testing"

func Benchmark_panic(t *testing.B) {
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		func() {
			defer func() {
				if e := recover(); e != nil {

				}
			}()
			panic(100)
		}()
	}
}
