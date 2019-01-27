package log

import (
	"runtime"
	"strings"
)

func GetRuntimeInfo(depList ...int) (file string, functionName string, line int) {
	var depth int
	if depList == nil {
		depth = 1
	} else {
		depth = depList[0]
	}
	function, file, line, _ := runtime.Caller(depth)
	functionName = runtime.FuncForPC(function).Name()
	return file, functionName, line
}

/*
formatter:
file: D:/dev/go/src/testing/testing.go
functionName: testing.tRunner
line: 8
*/
func GetRuntimeInfoShortFormat(depList ...int) (file string, functionName string, line int) {
	file, functionName, line = GetRuntimeInfo(depList...)
	fileIndex := strings.LastIndex(file, "/")
	funcIndex := strings.LastIndex(functionName, ".")
	return file[fileIndex+1:], functionName[funcIndex+1:], line
}
