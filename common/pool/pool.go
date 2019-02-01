package pool

import "sync"

const BufferSize = 4108

func init() {
	pool = &sync.Pool{
		New: func() interface{} {
			return make([]byte, BufferSize)
		},
	}
}

var pool *sync.Pool

func GetBuf() []byte {
	buf := pool.Get().([]byte)
	buf = buf[:cap(buf)]
	return buf
}

func PutBuf(buf []byte) {
	pool.Put(buf)
}

// type BytesPool struct {
// 	Size int
// 	*sync.Pool
// }

// func NewBytesPool(size int) *BytesPool {
// 	return &BytesPool{
// 		Size: size,
// 		Pool: &sync.Pool{
// 			New: func() interface{} {
// 				return make([]byte, size)
// 			},
// 		},
// 	}
// }

// func (this *BytesPool) Get() []byte {
// 	buf := this.Pool.Get().([]byte)
// 	buf = buf[:cap(buf)]
// 	return buf
// }

// func (this *BytesPool) Put(buf []byte) {
// 	this.Pool.Put(buf)
// }
