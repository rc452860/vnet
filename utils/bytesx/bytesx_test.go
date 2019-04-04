package bytesx

import (
	"fmt"
	"testing"
)

func ExampleContactSlice() {
	rst := ContactSlice([]byte{0x01, 0x02}, []byte{0x03, 0x04})
	fmt.Println(rst)
	//Output:
	//[1 2 3 4]
}

func BenchmarkContactSlice(t *testing.B) {
	for i := 0; i < t.N; i++ {
		ContactSlice([]byte{0x01, 0x02}, []byte{0x03, 0x04})
	}
}
