package httpd

import "time"

type Config struct {
	SessionStore         Store
	SessionKeepAliveTime time.Duration
}

func (c *Config) fill() {
	if c.SessionStore == nil {
		c.SessionStore = &defaultStore{}
	}
}
