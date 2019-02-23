package goroutine

import "log"

func Protect(g func()) {
	defer func() {
		if x := recover(); x != nil {
			log.Printf("run time panic: %v", x)
		}
	}()
	g()
}
