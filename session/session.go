package session

type Session interface {
	New(key string) Session
	ID() string
	Get(key string) any
	Set(key string, value any)
	UpdateExpire()
	IsExpired() bool
}
