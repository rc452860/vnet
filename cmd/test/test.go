package main

import (
	"fmt"
)

type Test struct {
	Name string
}

func main() {
	a := Test{}
	fmt.Print(a.(type))
}
