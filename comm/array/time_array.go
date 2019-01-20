package array

import (
	"context"
	"reflect"
	"sync"
	"time"
)

type TimeArrayElement struct {
	Time    time.Time
	Element interface{}
}

type TimeArray struct {
	array     []TimeArrayElement
	Size      int
	expire    time.Duration
	context   context.Context
	cancel    context.CancelFunc
	AutoClear bool
	sync.RWMutex
}

func NewTimeArray(expire time.Duration, autoClear bool) TimeArray {
	result := TimeArray{
		expire:    expire,
		AutoClear: autoClear,
	}
	if autoClear {
		ctx, cancel := context.WithCancel(context.Background())
		result.context = ctx
		result.cancel = cancel
		go result.autoClear()
	}
	return result
}

func (t *TimeArray) Close() {
	if t.AutoClear && t.cancel != nil {
		t.cancel()
	}
}

func (t *TimeArray) autoClear() {
	tick := time.Tick(t.expire)
	for {
		select {
		case <-t.context.Done():
			return
		case <-tick:
		}
		t.Clear()
	}
}

func (t *TimeArray) Add(data interface{}) {
	t.Lock()
	defer t.Unlock()

	t.array = append(t.array, TimeArrayElement{
		Time:    time.Now(),
		Element: data,
	})
	t.Size++
}

func (t *TimeArray) Remove(data interface{}) {
	for i, item := range t.array {
		if reflect.DeepEqual(data, item) {
			t.array = append(t.array[:i], t.array[i+1:]...)
			t.Size--
		}
	}
}

func (t *TimeArray) Clear() {
	now := time.Now()
	for i, item := range t.array {
		if now.Sub(item.Time) > t.expire {
			t.array = append(t.array[:i], t.array[i+1:]...)
			t.Size--
		}
	}
}
func (t *TimeArray) Range(callback func(i int, key interface{})) {
	for i, item := range t.array {
		callback(i, item.Element)
	}
}
