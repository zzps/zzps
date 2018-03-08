package zzps

import (
	"net/http"
	"strconv"
	"strings"
)

type Context struct {
	Request  *http.Request
	Response  http.ResponseWriter
}

func newContext() *Context {
	return new(Context)
}
func (context *Context) Reset(rw http.ResponseWriter, r *http.Request)  {
	context.Request = r
	context.Response = rw
}
//data model should implement this interface,
//then could bind parameters from request and
//response client with json data
type ZConverter interface {
	ToJson() string
	ToSctFromMap(map[string]string)
}
//bind to data model with request parameters
func (context *Context) BindTo(converter ZConverter){
	m := convertRequestValues(context.Request.Form)
	converter.ToSctFromMap(m)
}
//general func to response client whith a general response entity
func (context *Context) ResponseClient(re ResponseEntity )  {
	context.Response.Write(re.ToJson())
}
func convertRequestValues(m map[string][]string) map[string]string {
	mp := make(map[string]string)
	for k,v := range m{
		mp[k] = v[0]
	}
	return mp
}
const (
	Success = "success"
	Fail = "fail"
)
type ResponseEntity struct {
	Status string
	Message string
	Object string
	Count int
}

func (re ResponseEntity) ToJson() []byte  {
	b := strings.HasPrefix(re.Object, "{")
	s := "{\"status\":\""+re.Status+"\",\"message\":\""+re.Message+"\",\"count\":"+strconv.Itoa(re.Count)
	if b {
		s = s + ",\"object\":"+re.Object+"}"
	}else {
		s = s + ",\"object\":\""+re.Object+"\"}"
	}
	return []byte(s)
}

