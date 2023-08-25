package session

import "time"

type Session interface {
	Init(keepAliveTime time.Duration)
	Use(token string)
	Get(token, key string) any
	Set(token, key string, value any)
	Del(token, key string)
	KeepAlive(token string)
}
