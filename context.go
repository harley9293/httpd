package httpd

import (
	"encoding/json"
	log "github.com/harley9293/blotlog"
	"github.com/harley9293/httpd/session"
	"net/http"
	"reflect"
)

const IndexError = -1

type Context struct {
	Session session.Session

	r       *http.Request
	w       http.ResponseWriter
	service *Service

	fn reflect.Value

	index       int
	middlewares []MiddlewareFunc
}

func (c *Context) CreateSession(key string) {
	c.Session = c.service.NewSession(key)
}

func (c *Context) Error(status int, err error) {
	c.index = IndexError
	http.Error(c.w, err.Error(), status)
	log.Error("send|%s|%s|%d|%s|", c.r.Method, c.r.URL, status, err.Error())
}

func (c *Context) Json(st any) {
	c.w.Header().Set("Content-Type", "application/json")
	out, err := json.Marshal(st)
	if err != nil {
		c.Error(http.StatusInternalServerError, err)
		return
	}
	_, err = c.w.Write(out)
	if err != nil {
		c.Error(http.StatusInternalServerError, err)
		return
	}
	log.Info("send|%s|%s|200|%+v", c.r.Method, c.r.URL, st)
}

func (c *Context) String(s string) {
	c.w.Header().Set("Content-Type", "text/plain")
	_, err := c.w.Write([]byte(s))
	if err != nil {
		c.Error(http.StatusInternalServerError, err)
		return
	}
	log.Info("send|%s|%s|200|%s", c.r.Method, c.r.URL, s)
}

func (c *Context) Next() {
	if c.index == IndexError {
		return
	}

	c.index++
	if c.index <= len(c.middlewares) {
		c.middlewares[c.index-1](c)
	}
	c.CallHandler()
}

func (c *Context) CallHandler() {
	var params []reflect.Value
	if c.fn.Type().NumIn() > 1 {
		arg := reflect.New(c.fn.Type().In(0))
		contentType := c.r.Header.Get("Content-Type")
		switch contentType {
		case "application/json":
			err := json.NewDecoder(c.r.Body).Decode(arg.Interface())
			if err != nil {
				c.Error(http.StatusBadRequest, err)
				return
			}
		default:
			err := c.r.ParseForm()
			if err != nil {
				c.Error(http.StatusBadRequest, err)
				return
			}
			result := make(map[string]string)
			for key, values := range c.r.Form {
				if len(values) > 0 {
					result[key] = values[0] // only take the first value for each key
				}
			}
			formJson, _ := json.Marshal(result)
			err = json.Unmarshal(formJson, arg.Interface())
			if err != nil {
				c.Error(http.StatusBadRequest, err)
				return
			}
		}
		params = append(params, arg.Elem())
	}

	params = append(params, reflect.ValueOf(c))
	c.fn.Call(params)
}
