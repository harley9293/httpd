package _default

import "time"

type Store struct {
	sessionMap    map[string]map[string]any
	keepAliveTime time.Duration
}

func (s *Store) Init(keepAliveTime time.Duration) {
	s.sessionMap = make(map[string]map[string]any)
	s.keepAliveTime = keepAliveTime
}

func (s *Store) Use(token string) {
	ss, ok := s.sessionMap[token]
	if !ok || ss["__keep_alive_time"].(time.Time).Before(time.Now()) {
		s.sessionMap[token] = make(map[string]any)
	}
	s.KeepAlive(token)
}

func (s *Store) Get(token, key string) any {
	if _, ok := s.sessionMap[token]; !ok {
		return nil
	}

	return s.sessionMap[token][key]
}

func (s *Store) Set(token, key string, value any) {
	if _, ok := s.sessionMap[token]; !ok {
		return
	}

	s.sessionMap[token][key] = value
}

func (s *Store) Del(token, key string) {
	if _, ok := s.sessionMap[token]; ok {
		delete(s.sessionMap[token], key)
	}
}

func (s *Store) KeepAlive(token string) {
	ss, ok := s.sessionMap[token]
	if !ok {
		return
	}

	ss["__keep_alive_time"] = time.Now().Add(s.keepAliveTime)
}
