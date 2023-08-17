package httpd

import (
	log "github.com/harley9293/blotlog"
	"net/http"
)

type Service struct {
	srv *http.Server

	globalMiddlewares []MiddlewareFunc
	router            *router
}

func NewService() *Service {
	return &Service{
		router: newRouter(),
	}
}

func (m *Service) AddHandler(method, path string, f any, middleware ...MiddlewareFunc) {
	err := m.router.register(method, path, f, append(m.globalMiddlewares, middleware...)...)
	if err != nil {
		panic(err)
	}
}

func (m *Service) AddMiddleWare(f ...MiddlewareFunc) {
	m.globalMiddlewares = append(m.globalMiddlewares, f...)
}

func (m *Service) LinstenAndServe(address string) error {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Info("recv|%s|%s", r.Method, r.URL)
		ro := m.router.match(r.Method, r.URL.Path)
		if ro.empty {
			ro.middlewares = m.globalMiddlewares
		}
		c := newContext(r, w, ro)
		c.Next()
	})
	m.srv = &http.Server{Addr: address, Handler: serveMux}

	return m.srv.ListenAndServe()
}
