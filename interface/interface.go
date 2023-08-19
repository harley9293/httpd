package _interface

import "time"

type Store interface {
	Init(keepAliveTime time.Duration)
	Use(token string)
	Get(token, key string) any
	Set(token, key string, value any)
	Del(token, key string)
	KeepAlive(token string)
}

type Generator interface {
	Rand() string
}
