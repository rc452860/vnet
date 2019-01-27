package log

import (
	"os"
)

type LogFileWriter struct {
	FileName string
	File     *os.File
}

func LogFileWriterFactory(name string) *LogFileWriter {
	file, _ := OpenFile(name)
	log := &LogFileWriter{
		FileName: name,
		File:     file,
	}
	return log
}

func (this *LogFileWriter) Write(message string) {
	this.File.WriteString(message)
}

// remember after used need close file
func OpenFile(file string) (*os.File, error) {
	return os.OpenFile(file, os.O_APPEND|os.O_CREATE, 0644)
}
