package zzps

import "net/http"

type Context struct {
	Request  *http.Request
	Response  http.ResponseWriter
}

func newContext() *Context {
	return new(Context)
}
func (context *Context)Reset(rw http.ResponseWriter, r *http.Request)  {
	context.Request = r
	context.Response = rw
}
