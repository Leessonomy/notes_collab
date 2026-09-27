package app

import (
	"testing"

	"go.uber.org/fx"
)

func TestGraphIsComplete(t *testing.T) {
	if err := fx.ValidateApp(options()); err != nil {
		t.Fatal(err)
	}
}
