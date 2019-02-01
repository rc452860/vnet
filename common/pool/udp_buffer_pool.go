package pool

import "sync"

var udpPool *sync.Pool

const MAX_UDP_BUF_SIZE int = 65507

func init() {
	udpPool = &sync.Pool{
		New: func() interface{} {
			return make([]byte, MAX_UDP_BUF_SIZE)
		},
	}
}

func GetUdpBuf() []byte {
	buf := udpPool.Get().([]byte)
	buf = buf[:cap(buf)]
	return buf
}

func PutUdpBuf(buf []byte) {
	udpPool.Put(buf)
}
