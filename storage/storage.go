package storage

import (
	"sync"
	"time"
)

var speeds []*Speed = make([]*Speed, 4096)

type Record struct {
	LocalHost  string
	TargetHost string
	MiddleHost string
	Upload     int64
	Download   int64
	CreateTime time.Time
	Status     int
}

type Speed struct {
	UpSpeed   int64
	DownSpeed int64
	UpBytes   int64
	DownBytes int64
	downlock  *sync.Mutex
	uplock    *sync.Mutex
}

func NewSpeed() *Speed {
	tmp := &Speed{
		downlock: &sync.Mutex{},
		uplock:   &sync.Mutex{},
	}
	speeds = append(speeds, tmp)
	return tmp
}

func (this Speed) Upload(traffic int64) {
	this.uplock.Lock()
	this.UpBytes += traffic
	this.uplock.Unlock()
}

func (this Speed) Download(traffic int64) {
	this.downlock.Lock()
	this.DownBytes += traffic
	this.downlock.Unlock()
}
