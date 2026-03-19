package sctx

import (
	"testing"
)

func TestServiceContext(t *testing.T) {
	serviceCtx := NewServiceContext(
		WithName("service-context"),
	)

	if err := serviceCtx.Load(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := serviceCtx.Load(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := serviceCtx.Stop(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
