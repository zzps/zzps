package zzps

import (
	"net/http"
	"sync"
	"fmt"
)
//zzps路由器，mapping属性为真实路由规则
type ZRouter struct {
	mapping map[string] mappingEntity
	contextPool sync.Pool
}
//指定路由到的具体controller实体和该实体处理方法
type mappingEntity struct {
	function func(*Context)error
}
type HandlerFuc func(*Context)error
//创建并初始化ZRouter
func newZRouter() *ZRouter {
	router := new(ZRouter)
	router.mapping = make(map[string]mappingEntity,256)
	router.contextPool.New = func() interface{} {
		return newContext()
	}
	return router
}
//添加映射到路由器
func (router *ZRouter) addMapping(url string,handler func(*Context)error)  {
	if _,ok := router.mapping[url];ok{
		//打印存在相同映射，级别错误，此方法不做映射
		return
	}
	entity := mappingEntity{
		function: handler,
	}
	router.mapping[url] = entity
	return
}

//Implement http.Handler interface.
func (router *ZRouter)ServeHTTP(rw http.ResponseWriter, r *http.Request)  {
	context := router.contextPool.Get().(*Context)
	context.Reset(rw,r)
	defer router.contextPool.Put(context)
	urlPath := r.URL.Path
	if value,ok := router.mapping[urlPath];ok {
		value.function(context)
	}else {
		fmt.Println("404")
		//不存在请求方法返回404
	}
}

