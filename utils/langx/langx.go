package langx

import "log"

func Must(fn func() (interface{}, error)) interface{} {
	v, err := fn()
	if err != nil {
		log.Fatalln(err)
	}
	return v
}