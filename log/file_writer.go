package log

import (
	"os"

	"github.com/rc452860/vnet/utils"
)

type LogFileWriter struct {
	FileName string
	File     *os.File
}

func LogFileWriterFactory(name string) *LogFileWriter {
	file, _ := utils.OpenFile(name)
	log := &LogFileWriter{
		FileName: name,
		File:     file,
	}
	return log
}

func (this *LogFileWriter) Write(message string) {
	this.File.WriteString(message)
}
