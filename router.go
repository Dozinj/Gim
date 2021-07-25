package Gim

import (
	"log"
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node //为每种请求方式都新建一棵数
	handlers map[string]HandlerFunc
}

func newRouter()*router{
	return &router{roots:make(map[string]*node),
		handlers:make(map[string]HandlerFunc)}
}


//将pattern 拆分为part切片
func parsePattern(pattern string)[]string {
	vs := strings.Split(pattern, "/")
	parts := make([]string, 0)

	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

//路由注册绑定HandlerFunc方法
func (r *router)addRouter(method,pattern string,handler HandlerFunc) {
	log.Printf("Route %s - %s\n", method, pattern)
	parts := parsePattern(pattern)
	_, ok := r.roots[method]
	if !ok {
		//新建该请求方法trie数
		r.roots[method] = new(node)
	}
	r.roots[method].insert(pattern, parts, 0)

	//路由-方法表注册
	key := method + "-" + pattern
	r.handlers[key] = handler
}


//getRoute 函数中，还解析了:和*两种匹配符的参数，返回一个 map
//例如/p/go/doc匹配到/p/:lang/doc，解析结果为：{lang: "go"}，
///static/css/geek_tutu.css匹配到/static/*filepath，解析结果为{filepath: "css/geek_tutu.css"}。
func (r *router)getRoute(method string,path string)(*node,map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string, 0)

	//判断请求方法的存在性
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	//判断路径存在性
	n := root.search(searchParts, 0)
	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			//记录动态路由请求参数 用于从c.param获取
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}

			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}



//在调用匹配到的handler前，将解析出来的路由参数赋值给了c.Params
//在路由注册玩后，根据路由选择处理器方法
func (r *router)handler(c *Context) {
	node, params := r.getRoute(c.Method, c.Path)
	if node != nil {
		c.Params = params
		key := c.Method + "-" + c.Path
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
