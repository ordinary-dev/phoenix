package controllers

import (
	"testing"

	"github.com/ordinary-dev/phoenix/testutils"
)

// Check that all templates can be loaded.
func TestLoadTemplates(t *testing.T) {
	if err := testutils.ResetWorkingDir(); err != nil {
		t.Fatal(err)
	}

	if err := LoadTemplates(); err != nil {
		t.Fatal(err)
	}
}
