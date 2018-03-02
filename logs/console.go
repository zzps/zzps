package logs

import (
	"os"
	"sync"
)

type consoleWriter struct {
	writer *os.File
	sync.Mutex
}

func newConsoleWriter() *consoleWriter {
	return &consoleWriter{
		writer:os.Stdout,
	}
}
func (consoleWriter *consoleWriter) zWrite(level string,p []byte){
	consoleWriter.Lock()
	defer consoleWriter.Unlock()
	consoleWriter.writer.Write(p)
}