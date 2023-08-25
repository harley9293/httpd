package httpd

import (
	"bytes"
	"errors"
	"github.com/harley9293/httpd/generator"
	"github.com/harley9293/httpd/store"
	"net/http"
	"reflect"
	"testing"
	"time"
)

type MockResponseWriter struct {
	HeaderMap  http.Header
	Body       bytes.Buffer
	StatusCode int
}

func (rw *MockResponseWriter) Header() http.Header {
	return rw.HeaderMap
}

func (rw *MockResponseWriter) Write(bytes []byte) (int, error) {
	if string(bytes) == "{\"test\":\"error\"}" || string(bytes) == "error" {
		return 0, errors.New("error")
	}
	return rw.Body.Write(bytes)
}

func (rw *MockResponseWriter) WriteHeader(statusCode int) {
	rw.StatusCode = statusCode
}

func TestContext_UseSession(t *testing.T) {
	// mock request/response
	req, _ := http.NewRequest("GET", "http://localhost:8080/", nil)
	resp := &MockResponseWriter{HeaderMap: make(http.Header), Body: bytes.Buffer{}}

	// mock store
	s := &store.Default{}
	s.Init(time.Hour)

	// test UseSession first
	context := &Context{
		r:       req,
		w:       resp,
		routes:  nil,
		session: newSession(s),
		config:  &Config{SessionGenerator: &generator.Default{}},
	}
	context.UseSession()
	if context.session == nil {
		t.Error("context use session error")
	}
	context.session.Set("test", "test")

	// test UseSession second
	cookies := resp.Header()["Set-Cookie"]
	for _, cookie := range cookies {
		req.Header.Add("Cookie", cookie)
	}
	context.UseSession()
	if context.session.Get("test").(string) != "test" {
		t.Error("context use session error")
	}
}

func TestContext_Response(t *testing.T) {
	// mock request
	req, _ := http.NewRequest("GET", "http://localhost:8080/", nil)

	// test json success
	resp := &MockResponseWriter{HeaderMap: make(http.Header), Body: bytes.Buffer{}}
	context := &Context{
		r: req,
		w: resp,
	}
	context.Json(map[string]string{"test": "test"})
	if resp.Header().Get("Content-Type") != "application/json" {
		t.Error("context json error")
	}
	if resp.Body.String() != "{\"test\":\"test\"}" {
		t.Error("context json error")
	}

	// test json marshal error
	resp = &MockResponseWriter{HeaderMap: make(http.Header), Body: bytes.Buffer{}}
	context = &Context{
		r: req,
		w: resp,
	}
	context.Json(make(chan int))
	if resp.StatusCode != http.StatusInternalServerError {
		t.Error("context json error")
	}

	// test json write error
	resp = &MockResponseWriter{HeaderMap: make(http.Header), Body: bytes.Buffer{}}
	context = &Context{
		r: req,
		w: resp,
	}
	context.Json(map[string]string{"test": "error"})
	if resp.StatusCode != http.StatusInternalServerError {
		t.Error("context json error")
	}

	// test string success
	resp = &MockResponseWriter{HeaderMap: make(http.Header), Body: bytes.Buffer{}}
	context = &Context{
		r: req,
		w: resp,
	}
	context.String("test")
	if resp.Header().Get("Content-Type") != "text/plain" {
		t.Error("context string error")
	}
	if resp.Body.String() != "test" {
		t.Error("context string error")
	}

	// test string write error
	resp = &MockResponseWriter{HeaderMap: make(http.Header), Body: bytes.Buffer{}}
	context = &Context{
		r: req,
		w: resp,
	}
	context.String("error")
	if resp.StatusCode != http.StatusInternalServerError {
		t.Error("context string error")
	}
}

func TestContext_Next(t *testing.T) {
	// mock
	req, _ := http.NewRequest("GET", "http://localhost:8080/", nil)
	resp := &MockResponseWriter{HeaderMap: make(http.Header), Body: bytes.Buffer{}}
	f := func(c *Context) {}
	routes := &routes{"/", "GET", reflect.ValueOf(f), nil, false}
	context := &Context{
		r:      req,
		w:      resp,
		routes: routes,
	}

	// test next success
	context.Next()

	// test next index error
	context.index = IndexError
	context.Next()

	// test next middleware
	middlewareFunc := func(c *Context) {
		c.Next()
	}
	context.routes.middlewares = append(context.routes.middlewares, middlewareFunc)
	context.index = 0
	context.Next()
}

func TestContext_CallHandler(t *testing.T) {
	// test empty handler
	req, _ := http.NewRequest("GET", "http://localhost:8080/", nil)
	resp := &MockResponseWriter{HeaderMap: make(http.Header), Body: bytes.Buffer{}}
	ro := &routes{"/", "GET", reflect.Value{}, nil, true}
	context := &Context{
		r:      req,
		w:      resp,
		routes: ro,
	}
	context.callHandler()
	if resp.StatusCode != http.StatusNotFound {
		t.Error("context call handler error")
	}

	// test json param success
	type Test struct {
		Test  string `json:"test"`
		Test2 int    `json:"test2"`
	}
	req, _ = http.NewRequest("GET", "http://localhost:8080/", bytes.NewBuffer([]byte("{\"test\":\"hello\"}")))
	req.Header.Set("Content-Type", "application/json")
	resp = &MockResponseWriter{HeaderMap: make(http.Header), Body: bytes.Buffer{}}
	f := func(req *Test, c *Context) { c.String(req.Test) }
	ro = &routes{"/", "GET", reflect.ValueOf(f), nil, false}
	context = &Context{
		r:      req,
		w:      resp,
		routes: ro,
	}
	context.callHandler()
	if resp.Body.String() != "hello" {
		t.Error("context call handler error")
	}

	// test json param error
	req, _ = http.NewRequest("POST", "http://localhost:8080/", bytes.NewBuffer([]byte("{\"test\":1}")))
	req.Header.Set("Content-Type", "application/json")
	resp = &MockResponseWriter{HeaderMap: make(http.Header), Body: bytes.Buffer{}}
	ro = &routes{"/", "POST", reflect.ValueOf(f), nil, false}
	context = &Context{
		r:      req,
		w:      resp,
		routes: ro,
	}
	context.callHandler()
	if resp.StatusCode != http.StatusBadRequest {
		t.Error("context call handler error")
	}

	// test form param success
	req, _ = http.NewRequest("GET", "http://localhost:8080/?test=world", nil)
	resp = &MockResponseWriter{HeaderMap: make(http.Header), Body: bytes.Buffer{}}
	ro = &routes{"/", "GET", reflect.ValueOf(f), nil, false}
	context = &Context{
		r:      req,
		w:      resp,
		routes: ro,
	}
	context.callHandler()
	if resp.Body.String() != "world" {
		t.Error("context call handler error")
	}

	// test form param error
	req, _ = http.NewRequest("GET", "http://localhost:8080/?test=%GG", nil)
	resp = &MockResponseWriter{HeaderMap: make(http.Header), Body: bytes.Buffer{}}
	ro = &routes{"/", "GET", reflect.ValueOf(f), nil, false}
	context = &Context{
		r:      req,
		w:      resp,
		routes: ro,
	}
	context.callHandler()
	if resp.StatusCode != http.StatusBadRequest {
		t.Error("context call handler error")
	}

	// test form param error2
	req, _ = http.NewRequest("GET", "http://localhost:8080/?test2=hello", nil)
	resp = &MockResponseWriter{HeaderMap: make(http.Header), Body: bytes.Buffer{}}
	ro = &routes{"/", "GET", reflect.ValueOf(f), nil, false}
	context = &Context{
		r:      req,
		w:      resp,
		routes: ro,
	}
	context.callHandler()
	if resp.StatusCode != http.StatusBadRequest {
		t.Error("context call handler error")
	}
}
