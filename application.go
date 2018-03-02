package zzps

import (
	"net/http"
	"zzps/logs"
)

const (
	HTTP                  = "http"
	HTTPS                 = "https"
	TCP                   = "tcp"
)
//application
type application struct {
	protocol          string  //http|https
	network           string  //tcp|unix
	Handler           http.Handler
	server            *http.Server
	applicationConfig *applicationConfig
}
//
func newApplication() *application {
	application := new(application)
	application.applicationConfig = ApplicationConfig
	application.Handler = newZRouter()
	application.server = &http.Server{
		Handler : application.Handler,
		Addr :  application.applicationConfig.Address,
	}
	application.applicationConfig = ApplicationConfig
	application.protocol = application.applicationConfig.Protocol
	logs.StartZLogger()
	return application
}
//run the application
func (application *application)run() {
	application.serve()
}
//
func (application *application)serve() error{
	var error error
	switch application.protocol {
	case HTTP:
		error = application.server.ListenAndServe()
	case HTTPS:
		error = application.server.ListenAndServeTLS(application.applicationConfig.TlsCertFile,application.applicationConfig.TlsKeyFile)
	}
	return error
}

