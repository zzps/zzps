package zzps

import (
	"fmt"
	"strings"
)

const (
	version = "1.0.0"
	banner = "zzps web framework" + version + " URL: https://github.com/zzps/zzps"
)

type zzps struct {
	application *application
}
var z zzps
func Run() {
	z.application.run()
}
func Build() {
	z.application = newApplication()
}
func AddMapping(url string, handler func(*Context))  {
	var b bool
	router := z.application.Handler.(*zRouter)
	b = strings.HasPrefix(url, "/")
	if !b {
		url = "/"+url
	}
	b = strings.HasSuffix(url, router.urlSuffix)
	if !b {
		url = strings.Split(url,".")[0] + router.urlSuffix
	}
	z.application.Handler.(*zRouter).addMapping(url,handler)
}
func init()  {
	fmt.Println(banner[1:])
}
