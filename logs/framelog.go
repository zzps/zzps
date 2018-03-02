package logs

import (
	"sync"
	"os"
	"fmt"
)

const fameLoggerFileLocation  = logRootPath+"frameLogger.log"

type frameLogger struct {
	writer []zWriter
}
func newFrameLogger() *frameLogger {
	//fameLoggerFile, e := os.OpenFile(fameLoggerFileLocation, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	//if e != nil {
	//	log.Fatalln("fail to open fameLoggerFile !")
	//}
	fileWriter := new(frameFileWriter)
	var err error
	os.MkdirAll(defaultFileLogPath,0660)
	fileWriter.writer ,err = os.OpenFile(fameLoggerFileLocation, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		fmt.Println(err)
	}
	return &frameLogger{
		writer:[]zWriter{newConsoleWriter(),fileWriter},
	}
}

func (frameLogger *frameLogger) getZWriter() []zWriter{
	return frameLogger.writer
}

type frameFileWriter struct {
	writer *os.File
	sync.Mutex
}
func (frameFileWriter *frameFileWriter) zWrite(level string,p []byte){
	frameFileWriter.Lock()
	defer frameFileWriter.Unlock()
	frameFileWriter.writer.Write(p)
}