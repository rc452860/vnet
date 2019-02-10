package pool

import "sync"

const BufferSize = 4096

var (
	poolMap map[int]*sync.Pool
)

func init() {
	poolMap = make(map[int]*sync.Pool)
}

func GetBuf() []byte {
	pool := poolMap[BufferSize]
	if pool == nil {
		poolMap[BufferSize] = &sync.Pool{
			New: createAllocFunc(BufferSize),
		}
	}
	buf := poolMap[BufferSize].Get().([]byte)
	buf = buf[:cap(buf)]
	return buf
}

func GetBufBySize(size int) []byte {
	pool := poolMap[size]
	if pool == nil {
		poolMap[size] = &sync.Pool{
			New: createAllocFunc(size),
		}
	}
	buf := poolMap[size].Get().([]byte)
	buf = buf[:cap(buf)]
	return buf
}

func PutBuf(buf []byte) {
	poolMap[cap(buf)].Put(buf)
}

func createAllocFunc(size int) func() interface{} {
	return func() interface{} {
		return make([]byte, size)
	}
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
