package httpd

import "github.com/harley9293/httpd/interface"

type Session struct {
	store _interface.Store
	token string
}

func newSession(store _interface.Store) *Session {
	return &Session{store: store}
}

func (s *Session) Use(token string) {
	s.token = token
	s.store.Use(token)
}

func (s *Session) Get(key string) any {
	return s.store.Get(s.token, key)
}
func (s *Session) Set(key string, value any) {
	s.store.Set(s.token, key, value)
}
func (s *Session) Del(key string) {
	s.store.Del(s.token, key)
}
func (s *Session) KeepAlive() {
	s.store.KeepAlive(s.token)
}
