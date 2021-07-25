package Gim

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

//封装context 减少重写
type Context struct {
	Writer http.ResponseWriter
	Req *http.Request

	//request
	Path string
	Method string
	Params map[string]string   //value 为动态参数

	//response
	StatusCode int
}


func newContext(w http.ResponseWriter,r *http.Request)*Context{
	return &Context{
		Writer: w,
		Req: r,
		Path: r.URL.Path,
		Method: r.Method,
	}
}

func (c *Context)Param(key string)string{
	value,_:=c.Params[key]
	return value
}

func (c *Context)PostForm(key string)string{
	return c.Req.FormValue(key)
}

func (c *Context)Query(key string)string{
	return c.Req.URL.Query().Get(key)
}


func (c *Context)Status(code int){
	c.StatusCode=code
	if code==http.StatusOK {  // 解决 http: superfluous response.WriteHeader 控制台日志打印
		return
	}
	c.Writer.WriteHeader(code)
}

func (c *Context)SetHeader(key,value string){
	c.Writer.Header().Set(key,value)
}


//response
func (c *Context)String(code int,format string,value ...interface{}){
	c.SetHeader("Content-Type", "text/plain")   //w.WriteHeader 后 Set Header 是无效的,所以要提前
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format,value...)))
}

func (c *Context)Json(code int,obj interface{}){
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder:=json.NewEncoder(c.Writer)

	if err:=encoder.Encode(obj);err!=nil{
		http.Error(c.Writer,err.Error(),500)
	}
}

func (c *Context)Data(code int,data []byte){
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context)HTML(code int,html string){
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}

