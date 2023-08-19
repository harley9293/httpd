package _default

import "testing"

func TestGenerator_Rand(t *testing.T) {
	generator := &Generator{}
	if generator.Rand() == "" {
		t.Error("session generator rand error")
	}
}
