package zzps

const (
	ControllerUrl,MethodUrl = "ControllerUrl","MethodUrl"
)
type Controller struct {
	ControllerUrl string
	MethodUrl map[string]func(...interface{})
}
