package httpd

import "time"

type SessionManager interface {
	Init(keepAliveTime time.Duration)
	UseSession(token string) *Session
}

type Session struct {
	data map[string]any
}

func (s *Session) Get(key string) any {
	return s.data[key]
}

func (s *Session) Set(key string, value any) {
	s.data[key] = value
}

type defaultSession struct {
	sessionMap    map[string]*Session
	keepAliveTime time.Duration
}

func (s *defaultSession) Init(keepAliveTime time.Duration) {
	s.sessionMap = make(map[string]*Session)
	s.keepAliveTime = keepAliveTime
}

func (s *defaultSession) UseSession(token string) *Session {
	if session, ok := s.sessionMap[token]; ok {
		return session
	} else {
		s.sessionMap[token] = &Session{}
		return s.sessionMap[token]
	}
}
