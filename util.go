package twig

import (
	"net/http"
	"reflect"
	"runtime"
	"strings"
)

// 获取当前请求路径
func GetReqPath(r *http.Request) string {
	path := r.URL.RawPath

	if path == "" {
		path = r.URL.Path
	}

	return path
}

// 获取handler的名称
func HandlerName(h HandlerFunc) string {
	t := reflect.ValueOf(h).Type()
	if t.Kind() == reflect.Func {
		return runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
	}
	return t.String()
}

// HelloTwig! ~~
func HelloTwig(c Ctx) error {
	return c.Stringf(http.StatusOK, "Hello %s!", "Twig")
}

// 包装handler
func Enhance(handler HandlerFunc, m []MiddlewareFunc) HandlerFunc {
	if m == nil {
		return handler
	}

	h := handler
	for i := len(m) - 1; i >= 0; i-- {
		h = m[i](h)
	}
	return h

}

type Route struct {
	Name   string
	Path   string
	Method string
}

// 判断当前请求是否为AJAX
func IsAJAX(r *http.Request) bool {
	return strings.Contains(r.Header.Get(HeaderXRequestedWith), XMLHttpRequest)
}

// 设置关联关系
func assoc(i interface{}, t *Twig) {
	if linker, ok := i.(Assocer); ok {
		linker.Assoc(t)
	}
}

func GetPartner(id string, c Ctx) Partner {
	t := c.Twig()
	if p, ok := t.Partner(id); ok {
		return p
	}

	c.Logger().Panicf("Twig: Partner (%s) is not exist!", id)

	return nil
}
