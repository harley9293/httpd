package httpd

import (
	"github.com/harley9293/httpd/default"
	"testing"
)

func TestSession(t *testing.T) {
	store := &_default.Store{}
	store.Init(0)

	session := newSession(store)
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
