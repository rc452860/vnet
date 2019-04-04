package cache

import (
	"runtime"
	"sync"
	"time"
)

type LRU struct{
	*lrucache
}

type lrucache struct{
	mapping *sync.Map
	lruJanitor *lruJanitor
}

type lruelement struct {
	Expired time.Time
	Payload interface{}
	TTL time.Duration
}

func (l *lrucache) Put(key interface{}, payload interface{}, ttl time.Duration) {
	l.mapping.Store(key, &lruelement{
		Payload: payload,
		Expired: time.Now().Add(ttl),
		TTL: ttl,
	})
	l.len++
}

func (l *lrucache)  Get(key interface{}) interface{} {
	item, exist := l.mapping.Load(key)
	if !exist {
		return nil
	}
	elm := item.(*lruelement)
	// expired
	if time.Since(elm.Expired) > 0 {
		l.mapping.Delete(key)
		return nil
	}
	// lru strategy
	elm.Expired = time.Now().Add(elm.TTL)
	l.mapping.Store(key,elm)
	return elm.Payload
}

func (l *lrucache) Delete(key interface{}){
	l.mapping.Delete(key)
}

func (l *lrucache) IsExist(key interface{}) bool{
	_,exist := l.mapping.Load(key)
	return exist
}

func (l *lrucache) First() interface{}{
	var result interface{} = nil
	l.mapping.Range(func(key,value interface{}) bool{
		result = value
		return false
	})
	return result
}

func (l *lrucache) Len() int{
	var result int = 0
	l.mapping.Range(func(key,valye interface{})bool{
		result ++
		return true
	})
	return result
}

func (l *lrucache) Clean(){
	l.mapping.Range(func(k, v interface{}) bool {
		// key := k.(string)
		elm := v.(*lruelement)
		if time.Since(elm.Expired) > 0 {
			l.mapping.Delete(k)
		}
		return true
	})
}

type lruJanitor struct {
	interval time.Duration
	stop     chan struct{}
}

func (j *lruJanitor) process(c *lrucache) {
	ticker := time.NewTicker(j.interval)
	for {
		select {
		case <-ticker.C:
			c.Clean()
		case <-j.stop:
			ticker.Stop()
			return
		}
	}
}

func stopLruJanitor(c *Cache) {
	c.janitor.stop <- struct{}{}
}


// New return *Cache
func NewLruCache(interval time.Duration) *LRU {
	j := &lruJanitor{
		interval: interval,
		stop:     make(chan struct{}),
	}
	c := &lrucache{lruJanitor: j,len:0}
	go j.process(c)
	lru := &LRU{c}
	// this is very interesting,it worth be deep learning
	runtime.SetFinalizer(lru, stopJanitor)
	return lru
}
