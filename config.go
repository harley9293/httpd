package httpd

import (
	"github.com/harley9293/httpd/generator"
	"github.com/harley9293/httpd/orm"
	"github.com/harley9293/httpd/session"
	"time"
)

type Config struct {
	Session              session.Session
	SessionKeepAliveTime time.Duration
	SessionGenerator     generator.Generator
	ORM                  orm.ORM
}

func (c *Config) fill() {
	if c.Session == nil {
		c.Session = &session.Default{}
	}
	if c.SessionGenerator == nil {
		c.SessionGenerator = &generator.Default{}
	}
}
