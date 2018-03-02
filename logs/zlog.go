package logs


type zzLogger struct {
	writer []zWriter
}

func (zzLogger *zzLogger)getZWriter()[]zWriter{
	return zzLogger.writer
}
func newZzLogger() *zzLogger {
	zl := new(zzLogger)
	zWriter := logConfig.ZWriter
	for _,value := range zWriter{
		switch value {
		case Console:
			zl.writer = append(zl.writer, newConsoleWriter())
		case File:
			zl.writer = append(zl.writer, newFileWriter())
		}
	}
	return zl
}

