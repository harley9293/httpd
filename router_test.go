package httpd

import "testing"

func TestRouter_Register(t *testing.T) {
	r := newRouter()

	// invalid handler
	err := r.register("GET", "/", 1)
	if err == nil {
		t.Error("invalid handler should return error")
	}

	// not pointer param
	err = r.register("GET", "/", func(a int, c *Context) {})
	if err == nil {
		t.Error("not pointer param should return error")
	}

	// not *Context param
	err = r.register("GET", "/", func(a *int) {})
	if err == nil {
		t.Error("not *Context param should return error")
	}

	// same path and method
	err = r.register("GET", "/", func(c *Context) {})
	if err != nil {
		t.Error("first path and method should not return error")
	}

	err = r.register("GET", "/", func(c *Context) {})
	if err == nil {
		t.Error("same path and method should return error")
	}
}

func TestRouter_Match(t *testing.T) {
	r := newRouter()
	_ = r.register("GET", "/", func(c *Context) {})

	// found match
	route := r.match("GET", "/")
	if route.empty {
		t.Error("found match should not be empty")
	}

	// not found match
	route = r.match("GET", "/notfound")
	if !route.empty {
		t.Error("not found match should be empty")
	}
}
