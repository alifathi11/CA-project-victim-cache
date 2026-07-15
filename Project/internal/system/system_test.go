package system

import (
	"testing"
	"victimcacheproject/internal/config"
)

func TestDefaultSystemIsValid(t *testing.T) {
	if err := New(config.Default()).Validate(); err != nil {
		t.Fatal(err)
	}
}
