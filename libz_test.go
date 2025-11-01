package libz_test

import (
	"testing"

	"github.com/pgaskin/go-libz"
	_ "github.com/pgaskin/go-libz/embed"
)

func TestInitialize(t *testing.T) {
	if err := libz.Initialize(); err != nil {
		t.Errorf("failed to initialize: %v", err)
	}
}
