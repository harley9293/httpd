package httpd

import (
	"errors"
	"fmt"
	"reflect"
)

type routes struct {
	path        string
	method      string
	fn          reflect.Value
	middlewares []MiddlewareFunc

	empty bool
}

type router struct {
	data []*routes
}

func newRouter() *router {
	return &router{}
}

func handlerVerify(value reflect.Value) error {
	if value.Type().Kind() != reflect.Func {
		return errors.New("add handler type error")
	}

	for i := 0; i < value.Type().NumIn(); i++ {
		if value.Type().In(i).Kind() != reflect.Ptr {
			return errors.New("handler param must be pointer")
		}
	}

	if value.Type().In(value.Type().NumIn()-1) != reflect.TypeOf(&Context{}) {
		return errors.New("handler second param must be *Context")
	}

	return nil
}

func (r *router) register(method, path string, f any, middlewares ...MiddlewareFunc) error {
	for _, v := range r.data {
		if v.method == method && v.path == path {
			return errors.New(fmt.Sprintf("handler already exists,method:%s, path:%s", method, path))
		}
	}

	fn := reflect.ValueOf(f)
	err := handlerVerify(fn)
	if err != nil {
		return err
	}

	r.data = append(r.data, &routes{
		path:        path,
		method:      method,
		fn:          fn,
		middlewares: middlewares,
	})

	return nil
}

func (r *router) match(method, path string) *routes {
	for _, v := range r.data {
		if v.method == method && v.path == path {
			return v
		}
	}
	return &routes{empty: true}
}
