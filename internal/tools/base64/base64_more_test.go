package base64

import (
	"errors"
	"testing"

	apperr "github.com/nni-hotel/demo-mcp/internal/platform/errors"
)

func TestParseAlphabet(t *testing.T) {
	a, err := ParseAlphabet("standard")
	if err != nil || a != AlphabetStandard {
		t.Fatalf("standard: %v %v", a, err)
	}
	a, err = ParseAlphabet("url")
	if err != nil || a != AlphabetURL {
		t.Fatalf("url: %v %v", a, err)
	}
	_, err = ParseAlphabet("invalid")
	if err == nil {
		t.Fatal("expected error")
	}
	var ae *apperr.AppError
	if !errors.As(err, &ae) || ae.Code != apperr.CodeInvalidRequest {
		t.Fatalf("got %v", err)
	}
}

func TestEncodeDeterministic(t *testing.T) {
	req := Request{Input: "test", Alphabet: AlphabetStandard, Padding: true}
	r1, err := Encode(req)
	if err != nil {
		t.Fatal(err)
	}
	r2, err := Encode(req)
	if err != nil {
		t.Fatal(err)
	}
	if r1.Output != r2.Output {
		t.Fatalf("non-deterministic: %q vs %q", r1.Output, r2.Output)
	}
}
