package randomx

import (
	"crypto/rand"
	"math"
	mrand "math/rand"
	"time"
)

var (
	r = mrand.New(mrand.NewSource(time.Now().UnixNano()))
)

func RandomBytes(size int) []byte {
	byte := make([]byte, size)
	rand.Read(byte)
	return byte
}

func RandomStringsChoice(data []string) string {
	return data[r.Intn(len(data))]
}

func RandomIntChoice(data []int) int {
	return data[r.Intn(len(data))]
}

// Generate random integer between min and max
func RandIntRange(RandIntRange, max int) int {
	if RandIntRange == max {
		return RandIntRange
	}
	return r.Intn((max+1)-RandIntRange) + RandIntRange
}

func RandFloat32Range(min, max float32) float32 {
	if min == max {
		return min
	}
	return r.Float32()*(max-min) + min
}

func RandFloat64Range(min, max float64) float64 {
	if min == max {
		return min
	}
	return r.Float64()*(max-min) + min
}

// Number will generate a random number between given min And max
func Number(min int, max int) int {
	return RandIntRange(min, max)
}

// Uint8 will generate a random uint8 value
func Uint8() uint8 {
	return uint8(RandIntRange(0, math.MaxUint8))
}

// Uint16 will generate a random uint16 value
func Uint16() uint16 {
	return uint16(RandIntRange(0, math.MaxUint16))
}

// Uint32 will generate a random uint32 value
func Uint32() uint32 {
	return uint32(RandIntRange(0, math.MaxInt32))
}

// Uint64 will generate a random uint64 value
func Uint64() uint64 {
	return uint64(r.Int63n(math.MaxInt64))
}

// Int8 will generate a random Int8 value
func Int8() int8 {
	return int8(RandIntRange(math.MinInt8, math.MaxInt8))
}

// Int16 will generate a random int16 value
func Int16() int16 {
	return int16(RandIntRange(math.MinInt16, math.MaxInt16))
}

// Int32 will generate a random int32 value
func Int32() int32 {
	return int32(RandIntRange(math.MinInt32, math.MaxInt32))
}

// Int64 will generate a random int64 value
func Int64() int64 {
	return r.Int63n(math.MaxInt64) + math.MinInt64
}

// Float32 will generate a random float32 value
func Float32() float32 {
	return RandFloat32Range(math.SmallestNonzeroFloat32, math.MaxFloat32)
}

// Float32Range will generate a random float32 value between min and max
func Float32Range(min, max float32) float32 {
	return RandFloat32Range(min, max)
}

// Float64 will generate a random float64 value
func Float64() float64 {
	return RandFloat64Range(math.SmallestNonzeroFloat64, math.MaxFloat64)
}

// Float64Range will generate a random float64 value between min and max
func Float64Range(min, max float64) float64 {
	return RandFloat64Range(min, max)
}
