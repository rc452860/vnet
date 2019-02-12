package pool

import (
	"bytes"
	"fmt"
)

func ExampleGetBuf() {
	buf := GetBuf()
	fmt.Printf("len: %v\n", len(buf))
	fmt.Printf("cap: %v\n", cap(buf))

	buf2 := buf[1024:]
	fmt.Printf("len: %v\n", len(buf2))
	fmt.Printf("cap: %v\n", cap(buf2))

	bufReader := bytes.NewBuffer(buf2)
	bufReader.Reset()
	fmt.Println(bufReader.Len())
	fmt.Printf("len: %v\n", len(buf2))
	fmt.Printf("cap: %v\n", cap(buf2))
	//Output:
	//len: 4096
	//cap: 4096
	//len: 3072
	//cap: 3072
}
