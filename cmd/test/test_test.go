package main

import (
	"database/sql"
	"fmt"
)

func ExampleTest() {
	a := &sql.NullInt64{}
	a = &sql.NullInt64{
		Int64: 0,
		Valid: false,
	}
	fmt.Println(a)
	//Output:
}
