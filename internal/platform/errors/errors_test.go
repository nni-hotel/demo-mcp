package errors

import (
	"errors"
	"net/http"
	"testing"
)

func TestAppError(t *testing.T) {
	e := InvalidBase64("bad input")
	if e.Code != CodeInvalidBase64 {
		t.Fatalf("code: got %q", e.Code)
	}
	if e.Status != http.StatusBadRequest {
		t.Fatalf("status: got %d", e.Status)
	}
	if StatusFor(e) != http.StatusBadRequest {
		t.Fatalf("StatusFor: got %d", StatusFor(e))
	}
}

func TestStatusForUnknown(t *testing.T) {
	if StatusFor(errors.New("other")) != http.StatusInternalServerError {
		t.Fatal("expected 500 for unknown error")
	}
}

func TestPayloadTooLarge(t *testing.T) {
	e := PayloadTooLarge(1024)
	if e.Code != CodePayloadTooLarge {
		t.Fatalf("code: got %q", e.Code)
	}
	if e.Status != http.StatusRequestEntityTooLarge {
		t.Fatalf("status: got %d", e.Status)
	}
}
