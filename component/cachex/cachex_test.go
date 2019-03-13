package cachex

import (
	"fmt"
)

func ExampleBuildKey() {
	key, err := BuildKey(1, 13)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(key)
	//Output:
}
