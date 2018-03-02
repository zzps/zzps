package logs

import (
	"os"
	"encoding/json"
	"log"
	"io/ioutil"
	"strings"
	"strconv"
)

const (
	logConfigDir = "./config/"
	logConfigFile = "zLog.json"
	defaultFileName = "zzpsApp"
	defaultFileSuffix = ".log"
)
const (
	Console = "console"
	File = "file"
	Mail = "mail"
	Sms = "sms"
)
type zLogConfig struct {
	LevelString string `json:"levelString"`
	LevelInt int8 `json:"levelInt"`
	EnableZlog string `json:"enableZlog"`
	FileName            string `json:"fileName"`
	EnableDailyRotate   string   `json:"enableDailyRotate"`
	EnableLevelSeparate string   `json:"enableLevelSeparate"`
	Suffix              string `json:"suffix"`
	ZWriter []string `json:"zWriter"`
}
var logConfig = newZLogConfig()
func newZLogConfig()*zLogConfig  {
	logConfig := new(zLogConfig)
	logConfig.FileName = defaultFileName
	logConfig.Suffix = defaultFileSuffix
	logConfig.LevelInt = LevelInfo
	logConfig.LevelString = Info
	logConfig.EnableLevelSeparate = "true"
	logConfig.EnableDailyRotate = "true"
	logConfig.EnableZlog = "false"
	logConfig.ZWriter = []string{Console,File}
	var logConfigFileName = logConfigDir + logConfigFile
	var err error
	os.MkdirAll(logConfigDir,0660)
	_, err = os.Stat(logConfigFileName)
	if err != nil {
		//not exist,create file and use default set
		applicationConfigByte, err := json.MarshalIndent(logConfig,"","    ")
		if err != nil {
			log.Fatalln()
		}
		applicationConfigFile, err := os.OpenFile(logConfigFileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
		if err != nil {
			log.Fatalln()
		}
		_, err = applicationConfigFile.Write(applicationConfigByte)
		if err != nil {
			log.Fatalln()
		}
		return logConfig
	}
	bytes, err := ioutil.ReadFile(logConfigFileName)
	if err != nil {
		log.Fatalln()
	}
	zLogConfigTemp := new(zLogConfig)
	err = json.Unmarshal(bytes, zLogConfigTemp)
	if err != nil {
		log.Fatalln(err)
	}
	return replaceLogConfigDefaultPropertyUseCustom(logConfig,zLogConfigTemp)
}
func replaceLogConfigDefaultPropertyUseCustom(defaultConfig, customConfig *zLogConfig) *zLogConfig {
	var err error
	if customConfig.LevelString != "" {
		defaultConfig.LevelString = customConfig.LevelString
	}
	if customConfig.LevelInt != 0 {
		defaultConfig.LevelInt = customConfig.LevelInt
	}
	if customConfig.EnableZlog != ""{
		_, err = strconv.ParseBool(customConfig.EnableZlog)
		if err != nil {
			log.Fatalln("zLogConfig.enableZlog value parse bool failure")
		}
		defaultConfig.EnableZlog = customConfig.EnableZlog
	}
	if customConfig.Suffix != "" {
		defaultConfig.Suffix = customConfig.Suffix
	}
	if customConfig.FileName != ""{
		contains := strings.HasSuffix(customConfig.FileName, defaultConfig.Suffix)
		if contains {
			defaultConfig.FileName = strings.TrimSuffix(customConfig.FileName,defaultConfig.Suffix)
		}else {
			split := strings.Split(customConfig.FileName, ".")
			defaultConfig.FileName = split[0]
		}
	}
	if customConfig.EnableDailyRotate != "" {
		_, err = strconv.ParseBool(customConfig.EnableDailyRotate)
		if err != nil {
			log.Fatalln("zLogConfig.enableDailyRotate value parse bool failure")
		}
		defaultConfig.EnableDailyRotate = customConfig.EnableDailyRotate
	}
	if customConfig.EnableLevelSeparate != "" {
		_, err = strconv.ParseBool(customConfig.EnableLevelSeparate)
		if err != nil {
			log.Fatalln("zLogConfig.enableLevelSeparate value parse bool failure")
		}
		defaultConfig.EnableLevelSeparate = customConfig.EnableLevelSeparate
	}
	return defaultConfig
}