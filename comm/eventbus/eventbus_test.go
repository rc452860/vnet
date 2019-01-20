package eventbus

import (
	"testing"

	"github.com/asaskevich/EventBus"
)

func Benchmark_eventbus(t *testing.B) {
	evbus := EventBus.New()
	evbus.SubscribeAsync("abc", func(arg string) {
		// t.Log(arg)
	}, false)

	go func() {
		for i := 0; i < t.N; i++ {
			evbus.Publish("abc", "aa")
		}
	}()
	go func() {
		for i := 0; i < t.N; i++ {
			evbus.Publish("abc", "aa")
		}
	}()
	for i := 0; i < t.N; i++ {
		evbus.Publish("abc", "aa")
	}

}

func Benchmark_Channel(t *testing.B) {
	ch := make(chan string, 128)
	go func() {
		for {
			<-ch
		}
	}()
	go func() {
		for i := 0; i < t.N; i++ {
			ch <- "aa"
		}
	}()
	go func() {
		for i := 0; i < t.N; i++ {
			ch <- "aa"
		}
	}()
	for i := 0; i < t.N; i++ {
		ch <- "aa"
	}
}
