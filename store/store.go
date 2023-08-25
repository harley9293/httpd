package store

import "time"

type Default struct {
	sessionMap    map[string]map[string]any
	keepAliveTime time.Duration
}

func (s *Default) Init(keepAliveTime time.Duration) {
	s.sessionMap = make(map[string]map[string]any)
	s.keepAliveTime = keepAliveTime
}

func (s *Default) Use(token string) {
	ss, ok := s.sessionMap[token]
	if !ok || ss["__keep_alive_time"].(time.Time).Before(time.Now()) {
		s.sessionMap[token] = make(map[string]any)
	}
	s.KeepAlive(token)
}

func (s *Default) Get(token, key string) any {
	if _, ok := s.sessionMap[token]; !ok {
		return nil
	}

	return s.sessionMap[token][key]
}

func (s *Default) Set(token, key string, value any) {
	if _, ok := s.sessionMap[token]; !ok {
		return
	}

	s.sessionMap[token][key] = value
}

func (s *Default) Del(token, key string) {
	if _, ok := s.sessionMap[token]; ok {
		delete(s.sessionMap[token], key)
	}
}

func (s *Default) KeepAlive(token string) {
	ss, ok := s.sessionMap[token]
	if !ok {
		return
	}

	ss["__keep_alive_time"] = time.Now().Add(s.keepAliveTime)
}
