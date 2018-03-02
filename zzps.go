package zzps

import "fmt"

const (
	version = "1.0.0"
	banner = "z web framework" + version + "  URL --------n"
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
func AddMapping(url string, handler func(*Context)error)  {
	z.application.Handler.(*ZRouter).addMapping(url,handler)
}
func init()  {
	fmt.Println(banner[1:])
}
