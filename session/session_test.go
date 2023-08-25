package session

import "testing"

func TestDefaultStore(t *testing.T) {
	store := &Default{}
	store.Init(0)

	// before Use
	if store.Get("test", "test") != nil {
		t.Error("session get error")
	}

	store.Set("test", "test", "test")
	if store.Get("test", "test") != nil {
		t.Error("session set error")
	}

	store.Del("test", "test")
	store.KeepAlive("test")

	// after Use
	store.Use("test")
	store.Set("test", "test", "test")
	if store.Get("test", "test").(string) != "test" {
		t.Error("session set error")
	}
	store.Del("test", "test")
	if store.Get("test", "test") != nil {
		t.Error("session del error")
	}
	store.Use("test")
	store.KeepAlive("test")
	if store.Get("test", "__keep_alive_time") == nil {
		t.Error("session keep alive error")
	}
}
