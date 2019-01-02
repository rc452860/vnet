package log

import (
	"os"
)

type TerminalWriter struct {
	File *os.File
}

func LogTerminalWriterFactory() *TerminalWriter {
	log := &TerminalWriter{
		File: os.Stdout,
	}
	return log
}

func (this *TerminalWriter) Write(message string) {
	this.File.WriteString(message)
}
