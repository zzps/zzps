package zzps

import (
	"net/http"
	"sync"
	"zzps/logs"
	"strings"
	"strconv"
)
//zzps router,hold the actual relation with url and handler
type zRouter struct {
	mapping map[string] mappingEntity
	contextPool sync.Pool
	enableUrlSuffix bool
	urlSuffix string
	enableStaticSupport bool
	staticDir string
	staticUrlPrefix string
}
//
type mappingEntity struct {
	function func(*Context)
}
type HandlerFunc func(*Context)
func newZRouter() *zRouter {
	router := new(zRouter)
	router.mapping = make(map[string]mappingEntity,256)
	router.contextPool.New = func() interface{} {
		return newContext()
	}
	var enable bool
	enable, _ = strconv.ParseBool(ApplicationConfig.EnableUrlSuffix)
	if enable {
		router.enableUrlSuffix = enable
		router.urlSuffix = ApplicationConfig.UrlSuffix
	}
	enable, _ = strconv.ParseBool(ApplicationConfig.EnableStaticSupport)
	if enable {
		router.enableStaticSupport = enable
		router.staticDir = ApplicationConfig.StaticDir
		router.staticUrlPrefix = ApplicationConfig.StaticUrlPrefix
	}
	return router
}
//add mapping to router
func (router *zRouter) addMapping(url string,handler func(*Context)) bool {
	if _,ok := router.mapping[url];ok{
		logs.ZLogger.Warn("the url:{} is exist,addMapping failure",url)
		return false
	}
	entity := mappingEntity{
		function: handler,
	}
	router.mapping[url] = entity
	return true
}

//Implement http.Handler interface.
func (router *zRouter)ServeHTTP(rw http.ResponseWriter, r *http.Request)  {
	context := router.contextPool.Get().(*Context)
	context.Reset(rw,r)
	defer router.contextPool.Put(context)
	urlPath := r.URL.Path
	if router.enableStaticSupport {
		if strings.HasPrefix(urlPath,router.staticUrlPrefix){
			filePath := router.staticDir + strings.TrimPrefix(urlPath,router.staticUrlPrefix)
			http.ServeFile(rw,r,filePath)
			return
		}
	}
	if router.enableUrlSuffix {
		if b := strings.HasSuffix(urlPath, router.urlSuffix);!b {
			return
		}
	}
	if value,ok := router.mapping[urlPath];ok {
		value.function(context)
	}else {
		rw.Write(ResponseEntity{Status:Fail,Message:"404"}.ToJson())
	}
}

