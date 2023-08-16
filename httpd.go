package httpd

import (
	"errors"
	log "github.com/harley9293/blotlog"
	"github.com/harley9293/httpd/session"
	"net/http"
	"time"
)

type MiddlewareFunc func(*Context)

type Service struct {
	srv *http.Server

	globalMiddlewares []MiddlewareFunc
	router            *router

	sessionMap  map[string]session.Session
	baseSession session.Session
}

func NewService() *Service {
	return &Service{
		router:      newRouter(),
		sessionMap:  make(map[string]session.Session),
		baseSession: &session.Default{CfgExpireTime: 24 * time.Hour},
	}
}

func (m *Service) AddHandler(method, path string, f any, middleware ...MiddlewareFunc) {
	err := m.router.add(method, path, f, middleware...)
	if err != nil {
		panic(err)
	}
}

func (m *Service) AddGlobalMiddleWare(f ...MiddlewareFunc) {
	m.globalMiddlewares = append(m.globalMiddlewares, f...)
}

func (m *Service) UseSession(session session.Session) {
	m.baseSession = session
}

func (m *Service) GetSession(id string) session.Session {
	if s, ok := m.sessionMap[id]; ok {
		if s.IsExpired() {
			delete(m.sessionMap, id)
			return nil
		} else {
			return s
		}
	} else {
		return nil
	}
}

func (m *Service) NewSession(id string) session.Session {
	s := m.baseSession.New(id)
	m.sessionMap[s.ID()] = s
	return s
}

func (m *Service) LinstenAndServe(address string) error {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Info("recv|%s|%s", r.Method, r.URL)
		c := &Context{
			r:           r,
			w:           w,
			service:     m,
			middlewares: m.globalMiddlewares,
		}
		if r.Method != http.MethodOptions {
			ro := m.router.get(r.Method, r.URL.Path)
			if ro == nil {
				c.Error(http.StatusNotFound, errors.New(http.StatusText(http.StatusNotFound)))
				return
			}
			if len(ro.middlewares) > 0 {
				c.middlewares = append(c.middlewares, ro.middlewares...)
			}
			c.fn = ro.fn
		}
		c.Next()
	})
	m.srv = &http.Server{Addr: address, Handler: serveMux}

	return m.srv.ListenAndServe()
}
