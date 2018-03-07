package zzps

import (
	"io/ioutil"
	"encoding/json"
	"log"
	"os"
	"strconv"
)

const (
	applicationConfigDir  = "./config/"
	applicationConfigFile = "application.json"
	defaultAddress        = "0.0.0.0:8080"
	defaultProtocol       = "http"
	defaultUrlSuffix      = ".zson"
	defaultStaticDir      = "./static"
	defaultStaticUrlPrefix = "/static"
)

type applicationConfig struct {
	Address string `json:"address"`
	Protocol string `json:"protocol"`
	TlsCertFile string `json:"tlsCertFile"`
	TlsKeyFile string `json:"tlsKeyFile"`
	EnableUrlSuffix string `json:"enableUrlSuffix"`
	UrlSuffix string `json:"urlSuffix"`
	EnableStaticSupport string `json:"enableStaticSupport"`
	StaticDir string `json:"staticDir"`
	StaticUrlPrefix string `json:"staticUrlPrefix"`
}

var ApplicationConfig  = newApplicationConfig()
func newApplicationConfig() *applicationConfig {
	//initialize and set default property
	appConfig := new(applicationConfig)
	appConfig.Address = defaultAddress
	appConfig.Protocol = defaultProtocol
	appConfig.EnableUrlSuffix = "true"
	appConfig.UrlSuffix = defaultUrlSuffix
	appConfig.EnableStaticSupport = "true"
	appConfig.StaticDir = defaultStaticDir
	appConfig.StaticUrlPrefix = defaultStaticUrlPrefix
	var applicationFileName = applicationConfigDir + applicationConfigFile
	var err error
	os.MkdirAll(applicationConfigDir,0660)
	_, err = os.Stat(applicationFileName)
	if err != nil {
		//not exist,create file and use default set
		applicationConfigByte, err := json.MarshalIndent(appConfig,"","    ")
		if err != nil {
			log.Fatalln()
		}
		applicationConfigFile, err := os.OpenFile(applicationFileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
		if err != nil {
			log.Fatalln()
		}
		_, err = applicationConfigFile.Write(applicationConfigByte)
		if err != nil {
			log.Fatalln()
		}
		return appConfig
	}
	bytes, err := ioutil.ReadFile(applicationFileName)
	if err != nil {
		log.Fatalln()
	}
	applicationConfigTemp := new(applicationConfig)
	err = json.Unmarshal(bytes, applicationConfigTemp)
	if err != nil {
		log.Fatalln(err)
	}
	return replaceApplicationConfigDefaultPropertyUseCustom(appConfig,applicationConfigTemp)
}

func replaceApplicationConfigDefaultPropertyUseCustom(defaultConfig, customConfig *applicationConfig) *applicationConfig {
	var error error
	if customConfig.Address != "" {
		defaultConfig.Address = customConfig.Address
	}
	if customConfig.EnableUrlSuffix != "" {
		_,error = strconv.ParseBool(customConfig.EnableUrlSuffix)
		if error == nil {
			defaultConfig.EnableUrlSuffix = customConfig.EnableUrlSuffix
		}
	}
	if customConfig.UrlSuffix != "" {
		defaultConfig.UrlSuffix = customConfig.UrlSuffix
	}
	if customConfig.EnableStaticSupport != "" {
		_,error = strconv.ParseBool(customConfig.EnableStaticSupport)
		if error == nil {
			defaultConfig.EnableStaticSupport = customConfig.EnableStaticSupport
		}
	}
	if customConfig.StaticDir != "" {
		defaultConfig.StaticDir = customConfig.StaticDir
	}
	if customConfig.StaticUrlPrefix != "" {
		defaultConfig.StaticUrlPrefix = customConfig.StaticUrlPrefix
	}
	return defaultConfig
}