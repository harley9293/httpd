package httpd

import (
	"github.com/harley9293/httpd/generator"
	"github.com/harley9293/httpd/orm"
	"github.com/harley9293/httpd/store"
	"time"
)

type Config struct {
	SessionStore         store.Store
	SessionKeepAliveTime time.Duration
	SessionGenerator     generator.Generator
	ORM                  orm.ORM
}

func (c *Config) fill() {
	if c.SessionStore == nil {
		c.SessionStore = &store.Default{}
	}
	if c.SessionGenerator == nil {
		c.SessionGenerator = &generator.Default{}
	}
}
