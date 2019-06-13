package langx

import (
	"fmt"
	"testing"
)

func TestInArray(t *testing.T){
	r,_ := InArray("a",[]string{"a","b"})
	if !r {
		t.Fatal("r not exist")
	}
}

func ExampleFirstResult() {
	fmt.Println(FirstResult(Abc))

	// Output:
}
func Abc() (string,string){
	return "a","b"
}