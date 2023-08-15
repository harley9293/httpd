package session

import (
	"fmt"
	"time"
)

type Default struct {
	id         string
	data       map[string]any
	expireTime time.Time

	CfgExpireTime time.Duration
}

func (s *Default) New(key string) Session {
	return &Default{
		id:         fmt.Sprintf("%s%d", key, time.Now().UnixNano()),
		data:       make(map[string]any),
		expireTime: time.Now().Add(s.CfgExpireTime),

		CfgExpireTime: s.CfgExpireTime,
	}
}

func (s *Default) ID() string {
	return s.id
}

func (s *Default) Get(key string) any {
	return s.data[key]
}

func (s *Default) Set(key string, value any) {
	s.data[key] = value
}

func (s *Default) UpdateExpire() {
	if s.CfgExpireTime == 0 {
		return
	} else {
		s.expireTime = time.Now().Add(s.CfgExpireTime)
	}
}

func (s *Default) IsExpired() bool {
	if s.CfgExpireTime == 0 {
		return false
	} else {
		return s.expireTime.Before(time.Now())
	}
}
