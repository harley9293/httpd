package httpd

import (
	"net/http"
	"testing"
	"time"
)

func TestService(t *testing.T) {
	s := NewService(&Config{})

	s.AddMiddleWare(func(c *Context) {
		c.Next()
	})

	// test add handler success
	f := func(c *Context) {}
	err := s.AddHandler("GET", "/", f)
	if err != nil {
		t.Error("service add handler error")
	}

	// test add handler error
	err = s.AddHandler("GET", "/", f)
	if err == nil {
		t.Error("service add handler error")
	}

	go func() {
		_ = s.LinstenAndServe(":8080")
	}()
	time.Sleep(100 * time.Millisecond)

	req, _ := http.NewRequest("GET", "http://localhost:8080/", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error("service linsten and serve error")
	}
	if resp.StatusCode != http.StatusOK {
		t.Error("service linsten and serve error")
	}

	req, _ = http.NewRequest("GET", "http://localhost:8080/test", nil)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Error("service linsten and serve error")
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Error("service linsten and serve error")
	}
}

func TestService_Close(t *testing.T) {
	s := NewService(&Config{})
	go func() {
		_ = s.LinstenAndServe(":8080")
	}()
	time.Sleep(100 * time.Millisecond)
	_ = s.Close()
}
