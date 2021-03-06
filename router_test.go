package Gim

import (
	"fmt"
	"reflect"
	"testing"
)


func newTestRouter() *router {
	r := newRouter()
	r.addRouter("GET", "/", nil)
	r.addRouter("GET", "/hello/:name", nil)
	r.addRouter("GET", "/hello/b/c", nil)
	r.addRouter("GET", "/hi/:name", nil)
	r.addRouter("GET", "/assets/*filepath", nil)
	return r
}

func TestParsePattern(t *testing.T) {
	ok := reflect.DeepEqual(parsePattern("/p/:name"), []string{"p", ":name"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*"), []string{"p", "*"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*name/*"), []string{"p", "*name"})
	if !ok {
		t.Fatal("test parsePattern failed")
	}
}

func TestGetRoute(t *testing.T) {
	r := newTestRouter()
	n, ps := r.getRoute("GET", "/hello/geek-tutu")
	if n == nil {
		t.Fatal("nil shouldn't be returned")
	}
	if n.pattern != "/hello/:name" {
		t.Fatal("should match /hello/:name")
	}
	if ps["name"] != "geek-tutu" {
		t.Fatal("name should be equal to 'geek-tutu'")
	}
	fmt.Printf("matched path: %s, params['name']: %s\n", n.pattern, ps["name"])


	n, ps = r.getRoute("GET", "/assets/file1/file2")
	if n.children != nil {
		t.Fatal("shouldn't have children")
	}
	if n.pattern != "/assets/*filepath" {
		t.Fatal("should match /assets/*filepath")
	}
	if ps["filepath"] != "file1/file2" {
		t.Fatal("name should be equal to 'file1/file2'")
	}
	fmt.Printf("matched path: %s, params['filepath']: %s\n", n.pattern, ps["filepath"])
}
