package Gim

import (
	"net/http"
)

type HandlerFunc func(ctx *Context)

type Engine struct {
	router *router
}


//创建Engine实列
func New()*Engine{
	return &Engine{router: newRouter()}
}

//添加静态路由
func (engine *Engine)addRouter(method,patten string,handler HandlerFunc){
	engine.router.addRouter(method,patten,handler)
}

//添加GET路由
func (engine *Engine)GET(patten string,handler HandlerFunc){
	engine.addRouter("GET",patten,handler)
}

//添加POST路由
func (engine *Engine)POST(patten string,handler HandlerFunc){
	engine.addRouter("POST",patten,handler)
}

func (engine *Engine)ServeHTTP(w http.ResponseWriter,r *http.Request){
	c:=newContext(w,r)
	engine.router.handler(c)
}


func (engine *Engine)Run(addr string)error{
	return http.ListenAndServe(addr,engine) //只要传入任何实现了 ServerHTTP 接口的实例，所有的HTTP请求，就都交给了该实例处理了
}


