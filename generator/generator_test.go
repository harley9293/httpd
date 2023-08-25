package generator

import "testing"

func TestGenerator_Rand(t *testing.T) {
	generator := &Default{}
	if generator.Rand() == "" {
		t.Error("session generator rand error")
	}
}
