package Gim

import (
	"net/http"
)

type HandlerFunc func(ctx *Context)

type RouterGroup struct { //路由分组
	prefix      string        //前缀
	middlewares []HandlerFunc //支持中间件
	parent      *RouterGroup  //父节点
	engine      *Engine       //所有路由组共享实列
}

type Engine struct {
	RouterGroup *RouterGroup //Engine拥有RouterGroup所有的能力。
	router      *router
	groups      []*RouterGroup
}

//创建Engine实列
func New() *Engine {
	//将和路由有关的函数，都交给RouterGroup实现
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

//添加静态路由
func (engine *Engine) addRoute(method, patten string, handler HandlerFunc) {
	engine.router.addRouter(method, patten, handler)
}

//添加GET路由
func (engine *Engine) GET(patten string, handler HandlerFunc) {
	engine.addRoute("GET", patten, handler)
}

//添加POST路由
func (engine *Engine) POST(patten string, handler HandlerFunc) {
	engine.addRoute("POST", patten, handler)
}

//路由分组
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) addRoute(method, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	group.engine.addRoute(method, pattern, handler)
}

func (group *RouterGroup) GET(patten string, handler HandlerFunc) {
	group.addRoute("GET", patten, handler)
}

func (group *RouterGroup) POST(patten string, handler HandlerFunc) {
	group.addRoute("POST", patten, handler)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := newContext(w, r)
	engine.router.handler(c)
}

func (engine *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, engine) //只要传入任何实现了 ServerHTTP 接口的实例，所有的HTTP请求，就都交给了该实例处理了
}
