package Gim

import (
	"log"
	"net/http"
)

type router struct {
	Handlers map[string]HandlerFunc
}


func newRouter()*router{
	return &router{make(map[string]HandlerFunc)}
}

func (r *router)addRouter(method,pattern string,handler HandlerFunc){
	log.Printf("Route %s - %s\n", method, pattern)
	key:=method+"-"+pattern
	r.Handlers[key]=handler
}

func (r *router)handler(c *Context){
	key:=c.Method+"-"+c.Path
	if handler,ok:=r.Handlers[key];ok{
		handler(c)
	}else{
		c.String(http.StatusNotFound,"404 NOT FOUND: %s\n", c.Path)
	}
}
