package gim

import (
	"fmt"
	"net/http"
)

type HandlerFunc func(http.ResponseWriter,*http.Request)

type Engine struct {
	router map[string]HandlerFunc
}


//创建Engine实列
func New()*Engine{
	return &Engine{router: make(map[string]HandlerFunc)}
}

//添加静态路由
func (engine *Engine)addRouter(method,patten string,handler HandlerFunc){
	key:=method+"-"+patten
	engine.router[key]=handler
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
	key:=r.Method+"-"+r.URL.Path
	if hanlder,ok:= engine.router[key];ok{
		hanlder(w,r)
	}else{
		fmt.Fprintf(w,"404 NOT FOUND %s\n",r.URL)
	}
}

func (engine *Engine)Run(addr string)error{
	return http.ListenAndServe(addr,engine) //只要传入任何实现了 ServerHTTP 接口的实例，所有的HTTP请求，就都交给了该实例处理了
}


