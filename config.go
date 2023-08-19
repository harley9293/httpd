package httpd

import (
	"github.com/harley9293/httpd/default"
	"github.com/harley9293/httpd/interface"
	"time"
)

type Config struct {
	SessionStore         _interface.Store
	SessionKeepAliveTime time.Duration
	SessionGenerator     _interface.Generator
}

func (c *Config) fill() {
	if c.SessionStore == nil {
		c.SessionStore = &_default.Store{}
	}
	if c.SessionGenerator == nil {
		c.SessionGenerator = &_default.Generator{}
	}
}
