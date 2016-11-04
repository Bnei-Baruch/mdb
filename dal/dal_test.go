package dal

import "testing"

func TestInit(t *testing.T) {
	if success := Init(); !success {
		t.Error("Expected true, got ", success)
	}
}
