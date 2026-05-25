package base64

import (
	"strings"
	"testing"
)

func TestEncodeDecodeStandard(t *testing.T) {
	req := Request{Input: "hello", Alphabet: AlphabetStandard, Padding: true}
	enc, err := Encode(req)
	if err != nil {
		t.Fatal(err)
	}
	if enc.Output != "aGVsbG8=" {
		t.Fatalf("encode: got %q", enc.Output)
	}
	dec, err := Decode(Request{Input: enc.Output, Alphabet: AlphabetStandard, Padding: true})
	if err != nil {
		t.Fatal(err)
	}
	if dec.Output != "hello" {
		t.Fatalf("decode: got %q", dec.Output)
	}
}

func TestEncodeURLNoPadding(t *testing.T) {
	enc, err := Encode(Request{Input: "a", Alphabet: AlphabetURL, Padding: false})
	if err != nil {
		t.Fatal(err)
	}
	if enc.Output != "YQ" {
		t.Fatalf("got %q", enc.Output)
	}
}

func TestDecodeInvalid(t *testing.T) {
	_, err := Decode(Request{Input: "!!!", Alphabet: AlphabetStandard, Padding: true})
	if err == nil {
		t.Fatal("expected error")
	}
}

func BenchmarkEncode1KB(b *testing.B) {
	input := strings.Repeat("a", 1024)
	req := Request{Input: input, Alphabet: AlphabetStandard, Padding: true}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := Encode(req); err != nil {
			b.Fatal(err)
		}
	}
}
