package zzps

import (
	"io/ioutil"
	"encoding/json"
	"log"
	"os"
)

const (
	applicationConfigDir  = "./config/"
	applicationConfigFile = "application.json"
	defaultAddress        = "0.0.0.0:8080"
	defaultProtocol       = "http"
)

type applicationConfig struct {
	Address string `json:"address"`
	Protocol string `json:"protocol"`
	TlsCertFile string `json:"tlsCertFile"`
	TlsKeyFile string `json:"tlsKeyFile"`
}

var ApplicationConfig  = newApplicationConfig()
func newApplicationConfig() *applicationConfig {
	//initialize and set default property
	appConfig := new(applicationConfig)
	appConfig.Address = defaultAddress
	appConfig.Protocol = defaultProtocol
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
	if customConfig.Address != "" {
		defaultConfig.Address = customConfig.Address
	}
	return defaultConfig
}