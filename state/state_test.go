package state

import (
	"os"
	"testing"
)

func TestNewState(t *testing.T) {
	_, err := NewState("./", "testStateDb.txt")
	defer os.Remove("testStateDb.txt")
	if err != nil {
		t.Error("error while creating new state instance ", err)
	}
}
