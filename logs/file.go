package logs

import (
	"sync"
	"time"
	"os"
	"strconv"
	"log"
)

const (
	logRootPath = "./logs/"
	defaultFileLogPath = logRootPath + "app/"
	defaultFileKey = "defaultFileKey"
	logFileConnectSymbol = "_"
)

type fileWriter struct {
	sync.Mutex
	fileName            string //without suffix
	//logPath             string
	enableDailyRotate   bool
	enableLevelSeparate bool
	suffix              string
	fileMap             map[string]*os.File
}

func newFileWriter() *fileWriter {
	fileWriter := new(fileWriter)
	fileWriter.suffix = logConfig.Suffix
	fileWriter.fileName = logConfig.FileName

	var boo bool
	boo, _ = strconv.ParseBool(logConfig.EnableDailyRotate)
	fileWriter.enableDailyRotate = boo
	boo, _ = strconv.ParseBool(logConfig.EnableLevelSeparate)
	fileWriter.enableLevelSeparate = boo

	//4 kinds may be encountered
	fileWriter.fileMap = make(map[string]*os.File)
	if fileWriter.enableLevelSeparate {
		out : for i := logConfig.LevelInt; i < LevelFatal; i++  {
			for k,v := range LevelMap {
				if i == v {
					filePath := defaultFileLogPath + fileWriter.fileName + logFileConnectSymbol + k +fileWriter.suffix
					file := createFile(filePath)
					fileWriter.fileMap[k] = file
					continue out
				}
			}
		}
		if fileWriter.enableDailyRotate {
			//start new routine to do rotate on application running
			go fileWriter.dailyRotate()
		}
	}else {
		filePath := defaultFileLogPath + fileWriter.fileName + fileWriter.suffix
		file := createFile(filePath)
		fileWriter.fileMap[defaultFileKey] = file
		if fileWriter.enableDailyRotate {
			//start new routine to do rotate on application running
			go fileWriter.dailyRotate()
		}
	}
	return fileWriter
}
func (fileWriter *fileWriter) zWrite(level string,p []byte) {
	fileWriter.Lock()
	defer fileWriter.Unlock()
	if fileWriter.enableLevelSeparate {
		fileWriter.fileMap[level].Write(p)
	}else {
		fileWriter.fileMap[defaultFileKey].Write(p)
	}
}
func (fileWriter *fileWriter) dailyRotate() {
	now := time.Now()
	y, m, d := now.Add(24 * time.Hour).Date()
	nextDay := time.Date(y, m, d, 0, 0, 0, 0, now.Location())
	duration := time.Duration(nextDay.UnixNano() - now.UnixNano() + 100)
	timer := time.NewTimer(duration)
	<-timer.C
	fileWriter.doDailyRotate()
	ticker := time.NewTicker(time.Hour * 24)
	<-ticker.C
	fileWriter.doDailyRotate()
}
func (fileWriter *fileWriter) doDailyRotate() {
	now := time.Now()
	var oldFilePath, newFilePath string
	if fileWriter.enableLevelSeparate {
		for k := range fileWriter.fileMap {
			oldFilePath = defaultFileLogPath + fileWriter.fileName + logFileConnectSymbol + k +fileWriter.suffix
			newFilePath = defaultFileLogPath + fileWriter.fileName + logFileConnectSymbol + k + strconv.Itoa(now.Year()) + "-" + strconv.Itoa(int(now.Month())) + "-" + strconv.Itoa(now.Day()) + fileWriter.suffix
			fileWriter.doFileRotate(oldFilePath,newFilePath,k)
		}
		return 
	}
	oldFilePath = defaultFileLogPath + fileWriter.fileName + fileWriter.suffix
	newFilePath = defaultFileLogPath + fileWriter.fileName + logFileConnectSymbol + strconv.Itoa(now.Year()) + "-" + strconv.Itoa(int(now.Month())) + "-" + strconv.Itoa(now.Day()) + fileWriter.suffix
	fileWriter.doFileRotate(oldFilePath,newFilePath,defaultFileKey)
}
//fullPath should be path + realName
func createFile(fullPath string) *os.File {
	file, e := os.OpenFile(fullPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	if e != nil {
		log.Fatalln("create log file failure")
		return nil
	}
	return file
}
func (fileWriter *fileWriter)doFileRotate(oldFilePath,newFilePath,fileKey string,)  {
	//oldFilePath = defaultFileLogPath + fileWriter.fileName + logFileConnectSymbol + k +fileWriter.suffix
	fileInfo, e := os.Lstat(oldFilePath)
	if e != nil {
		log.Fatalln("oldFilePath may be wrong")
		return
	}
	size := fileInfo.Size()
	if size == 0 {
		//表明过去一天此等级的日志没有写，那就跳过日志日期循环
		return
	}
	//newFilePath = defaultFileLogPath + fileWriter.fileName + logFileConnectSymbol + k + strconv.Itoa(now.Year()) + "-" + strconv.Itoa(int(now.Month())) + "-" + strconv.Itoa(now.Day()) + fileWriter.suffix
	//重命名前先关闭写
	fileWriter.Lock()
	fileWriter.fileMap[fileKey].Close()
	os.Rename(oldFilePath, newFilePath)
	//再创建一个oldFileLocation文件，并将fileWriter的map里的指针换掉
	newOldFile := createFile(oldFilePath)
	fileWriter.fileMap[fileKey] = newOldFile
	fileWriter.Unlock()
}
