package main

import (
	"fmt"

	"github.com/rc452860/vnet/ciphers"
)

type Test struct {
	Name string
}

func main() {
	for _, item := range ciphers.GetSupportCiphers() {
		fmt.Println(item)
	}
}
