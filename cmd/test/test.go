package main

import (
	"fmt"
	"time"
)

func main() {
	go func() {
		defer func() {
			if e := recover(); e != nil {
				fmt.Printf("error %v", e)
			}
		}()
		go func() {
			defer func() {
				if e := recover(); e != nil {
					fmt.Printf("error %v \n", e)
				}
			}()
			panic("this is error")
		}()
	}()
	time.Sleep(1 * time.Second)
	fmt.Println("aaa")
}
