package _default

import "testing"

func TestDefaultStore(t *testing.T) {
	store := &Store{}
	store.Init(0)

	// before Use
	if store.Get("test", "test") != nil {
		t.Error("store get error")
	}

	store.Set("test", "test", "test")
	if store.Get("test", "test") != nil {
		t.Error("store set error")
	}

	store.Del("test", "test")
	store.KeepAlive("test")

	// after Use
	store.Use("test")
	store.Set("test", "test", "test")
	if store.Get("test", "test").(string) != "test" {
		t.Error("store set error")
	}
	store.Del("test", "test")
	if store.Get("test", "test") != nil {
		t.Error("store del error")
	}
	store.Use("test")
	store.KeepAlive("test")
	if store.Get("test", "__keep_alive_time") == nil {
		t.Error("store keep alive error")
	}
}
