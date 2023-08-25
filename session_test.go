package httpd

import (
	"github.com/harley9293/httpd/store"
	"testing"
)

func TestSession(t *testing.T) {
	s := &store.Default{}
	s.Init(0)

	session := newSession(s)
	session.Use("test")
	session.Set("test", "test")
	if session.Get("test").(string) != "test" {
		t.Error("session set error")
	}
	session.Del("test")
	if session.Get("test") != nil {
		t.Error("session del error")
	}
	session.KeepAlive()
	if session.Get("__keep_alive_time") == nil {
		t.Error("session keep alive error")
	}
}
