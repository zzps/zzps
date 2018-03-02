package logs

import (
	"strings"
	"time"
	"runtime"
	"strconv"
)
//日志级别
const(
	LevelTrace = iota   //一般细粒度信息
	LevelDebug          //调试信息
	LevelInfo           //程序运行情况一般信息
	LevelWarn           //程序运行存在潜在问题需要关注
	LevelError          //运行出现错误事件但程序仍可继续运行
	LevelFatal          //出现非常严重事件必须终止程序运行
)
const (
	Trace = "TRACE"
	Debug = "DEBUG"
	Info = "INFO"
	Warn = "WARN"
	Error = "ERROR"
	Fatal = "FATAL"
)
const (
	placeholder = "{}"
	)
var LevelMap = map[string]int8{Trace:LevelTrace,Debug:LevelDebug,Info:LevelInfo,Warn:LevelWarn,Error:LevelError,Fatal:LevelFatal}
//需要实现此接口才能用于打印日志
type ZLog interface {
	Trace(format string,args ...interface{})
	Debug(format string,args ...interface{})
	Info(format string,args ...interface{})
	Warn(format string,args ...interface{})
	Error(format string,args ...interface{})
	Fatal(format string,args ...interface{})
}

type zWriter interface {
	zWrite(level string,b []byte)
}
type zLoggerI interface {
	getZWriter()[]zWriter
}
type zLogger struct {
	levelString string
	levelInt int8
	enableCallLogLocation bool
	zLoggerI zLoggerI
	enableZlog bool    //false:framelog;true:zlog
}
var ZLogger = func()*zLogger{
	logger := new(zLogger)
	logger.levelString = logConfig.LevelString
	logger.levelInt = logConfig.LevelInt
	enableZlog,_ := strconv.ParseBool(logConfig.EnableZlog)
	if enableZlog {
		logger.zLoggerI = newZzLogger()
	}else {
		logger.zLoggerI = newFrameLogger()
	}
	logger.enableCallLogLocation = true
	return logger
}()
//Set log level by string flag
//example ZLogger.SetLevelString(Error)
func (zLogger *zLogger)SetLevelString(levelString string) bool {
	v,ok := LevelMap[levelString]
	if !ok {
		ZLogger.Error("please check the input params levelString")
		return false
	}
	zLogger.levelString = levelString
	zLogger.levelInt = v
	return true
}
//Set log level by int flag
//example ZLogger.SetLevelString(LevelError)
func (zLogger *zLogger)SetLevelInt(levelInt int8) bool {
	var boo = false
	for k,v := range LevelMap{
		if levelInt == v {
			boo = true
			zLogger.levelString = k
			zLogger.levelInt = v
			break
		}
	}
	if !boo {
		ZLogger.Error("please check the input params levelInt")
		return false
	}
	return true
}
func (zLogger *zLogger)Trace(format string,args ...interface{}){
	if zLogger.levelInt > LevelTrace {
		return
	}
	buf := zLogger.packaging(Trace,format, args)
	for _, writer := range zLogger.zLoggerI.getZWriter() {
		writer.zWrite(Trace,buf)
	}
}
func (zLogger *zLogger)Debug(format string,args ...interface{}){
	if zLogger.levelInt > LevelDebug {
		return
	}
	buf := zLogger.packaging(Debug,format, args)
	for _, writer := range zLogger.zLoggerI.getZWriter() {
		writer.zWrite(Debug,buf)
	}
}
func (zLogger *zLogger)Info(format string,args ...interface{}){
	if zLogger.levelInt > LevelInfo {
		return
	}
	buf := zLogger.packaging(Info,format, args)
	for _, writer := range zLogger.zLoggerI.getZWriter() {
		writer.zWrite(Info,buf)
	}
}
func (zLogger *zLogger)Warn(format string,args ...interface{}){
	if zLogger.levelInt > LevelWarn {
		return
	}
	buf := zLogger.packaging(Warn,format, args)
	for _, writer := range zLogger.zLoggerI.getZWriter() {
		writer.zWrite(Warn,buf)
	}
}
func (zLogger *zLogger)Error(format string,args ...interface{}){
	if zLogger.levelInt > LevelError {
		return
	}
	buf := zLogger.packaging(Error,format, args)
	for _, writer := range zLogger.zLoggerI.getZWriter() {
		writer.zWrite(Error,buf)
	}
}
func (zLogger *zLogger)Fatal(format string,args ...interface{}){
	if zLogger.levelInt > LevelFatal {
		return
	}
	buf := zLogger.packaging(Fatal,format, args)
	for _, writer := range zLogger.zLoggerI.getZWriter() {
		writer.zWrite(Fatal,buf)
	}
}
func formatHeader(buf []byte,level string,file string,line int) []byte {
	now := time.Now()
	formatTime := now.Format("2006-01-02 15:04:05")
	buf = append(buf, level+":"...)
	buf = append(buf,formatTime+":"...)
	buf = append(buf,file+":"...)
	buf = append(buf,strconv.Itoa(line)+":"...)
	return buf
}
func (zLogger *zLogger)getFileLine() (string, int) {
	var file string
	var line int
	if zLogger.enableCallLogLocation {
		var ok bool
		_, file, line, ok = runtime.Caller(2)
		if !ok {
			file = "???"
			line = 0
		}
		//_, filename := path.Split(file)
		//format = "[" + filename + ":" + strconv.Itoa(line) + "] " + format
	}
	return file,line
}
func (zLogger *zLogger)packaging(level string,format string,args ...interface{}) []byte {
	format = handlePlaceholder(format, args)
	file, line := zLogger.getFileLine()
	buf := make([]byte,0,35)
	buf = formatHeader(buf,level, file, line)
	buf = append(append(buf, format...),'\n')
	return buf
}
func handlePlaceholder(format string, args ...interface{}) string {
	length := len(args)
	if length ==0 {
		return format
	}
	count := strings.Count(format, placeholder)
	if count != length{
		return format
	}
	for i := 0; i < length; i++ {
		s, ok := args[i].(string)
		if !ok {
			continue
		}
		format = strings.Replace(format,placeholder,s,1)
	}
	return format
}
func StartZLogger()  {
	switch logConfig.LevelInt {
	case LevelTrace:
		ZLogger.Trace("built application ...")
	case LevelDebug:
		ZLogger.Debug("built application ...")
	case LevelInfo:
		ZLogger.Info("built application ...")
	case LevelWarn:
		ZLogger.Warn("built application ...")
	case LevelError:
		ZLogger.Error("built application ...")
	}
}