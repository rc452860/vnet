package main

import (
	"fmt"
)

func main() {
	s := "a,c,b,"
	fmt.Printf(s[:len(s)-1])

}
