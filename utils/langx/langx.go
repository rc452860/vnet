package langx

import (
	"log"
	"reflect"
)

func Must(fn func() (interface{}, error)) interface{} {
	v, err := fn()
	if err != nil {
		log.Fatalln(err)
	}
	return v
}

// get one result from multi result function
// example:
// func TestFunc() (string,string){
//     return "a","b"
// }
//
// func ExampleFirstResult(){
// 		FirstResult(TestFunc)
// }
// Output:
// 	a
func FirstResult(function interface{}, args... interface{}) interface{}{
	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	if v := reflect.ValueOf(function); v.Kind() == reflect.Func {
		results :=  v.Call(inputs)
		if len(results) > 0{
			return results[0].Interface()
		}else{
			return nil
		}
	}else {
		return nil
	}
}

func InArray(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}
	return
}