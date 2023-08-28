package httpd

import (
	"context"
	log "github.com/harley9293/blotlog"
	"github.com/harley9293/httpd/session"
	"net/http"
)

type Service struct {
	srv    *http.Server
	config *Config

	globalMiddlewares []MiddlewareFunc
	router            *router
	session           session.Session
}

func NewService(config *Config) *Service {
	config.fill()
	s := &Service{}
	s.config = config
	s.router = newRouter()
	s.session = config.Session
	s.session.Init(config.SessionKeepAliveTime)

	s.init()

	return s
}

func (m *Service) init() {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Info("recv|%s|%s", r.Method, r.URL)
		ro := m.router.match(r.Method, r.URL.Path)
		if ro.empty {
			ro.middlewares = m.globalMiddlewares
		}
		c := &Context{
			r:       r,
			w:       w,
			routes:  ro,
			session: m.session,
			config:  m.config,
		}
		c.Next()
	})
	m.srv = &http.Server{Addr: m.config.Address, Handler: serveMux}
}

func (m *Service) AddHandler(method, path string, f any, middleware ...MiddlewareFunc) error {
	return m.router.register(method, path, f, append(m.globalMiddlewares, middleware...)...)
}

func (m *Service) AddMiddleWare(f ...MiddlewareFunc) {
	m.globalMiddlewares = append(m.globalMiddlewares, f...)
}

func (m *Service) LinstenAndServe() error {
	return m.srv.ListenAndServe()
}

func (m *Service) Close() error {
	return m.srv.Shutdown(context.Background())
}
