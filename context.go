package httpd

import (
	"encoding/json"
	"errors"
	log "github.com/harley9293/blotlog"
	"net/http"
	"reflect"
)

const IndexError = -1

type MiddlewareFunc func(*Context)

type Context struct {
	r       *http.Request
	w       http.ResponseWriter
	routes  *routes
	session *Session
	config  *Config

	index int
}

func (c *Context) UseSession() *Session {
	cookie, err := c.r.Cookie("token")
	var token string
	if err == nil && cookie != nil {
		token = cookie.Value
	} else {
		token = c.config.SessionGenerator.Rand()
	}
	c.session.Use(token)

	http.SetCookie(c.w, &http.Cookie{
		Name:  "token",
		Value: token,
	})

	return c.session
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
	if c.index <= len(c.routes.middlewares) {
		c.routes.middlewares[c.index-1](c)
	}
	c.callHandler()
}

func (c *Context) callHandler() {
	if c.routes.empty {
		c.Error(http.StatusNotFound, errors.New(http.StatusText(http.StatusNotFound)))
		return
	}

	var params []reflect.Value
	if c.routes.fn.Type().NumIn() > 1 {
		arg := reflect.New(c.routes.fn.Type().In(0))
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
	c.routes.fn.Call(params)
}
