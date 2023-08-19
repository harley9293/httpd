package httpd

import "testing"

func TestConfig(t *testing.T) {
	config := &Config{}
	config.fill()
	if config.SessionStore == nil {
		t.Error("config fill error")
	}
}
