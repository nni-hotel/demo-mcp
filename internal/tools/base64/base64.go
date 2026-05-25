package base64

import (
	"encoding/base64"
	"strings"
	"time"

	apperr "github.com/nni-hotel/demo-mcp/internal/platform/errors"
)

type Alphabet string

const (
	AlphabetStandard Alphabet = "standard"
	AlphabetURL      Alphabet = "url"
)

type Request struct {
	Input    string
	Alphabet Alphabet
	Padding  bool
}

type Result struct {
	Output      string
	InputBytes  int64
	OutputBytes int64
	DurationMs  float64
	Tool        string
}

func Encode(req Request) (Result, error) {
	start := time.Now()
	enc, err := encoder(req.Alphabet, req.Padding)
	if err != nil {
		return Result{}, err
	}
	input := []byte(req.Input)
	out := enc.EncodeToString(input)
	elapsed := time.Since(start)
	return Result{
		Output:      out,
		InputBytes:  int64(len(input)),
		OutputBytes: int64(len(out)),
		DurationMs:  float64(elapsed.Microseconds()) / 1000.0,
		Tool:        "base64.encode",
	}, nil
}

func Decode(req Request) (Result, error) {
	start := time.Now()
	enc, err := encoder(req.Alphabet, req.Padding)
	if err != nil {
		return Result{}, err
	}
	input := req.Input
	decoded, err := enc.DecodeString(input)
	if err != nil {
		return Result{}, apperr.InvalidBase64("input is not valid base64")
	}
	out := string(decoded)
	elapsed := time.Since(start)
	return Result{
		Output:      out,
		InputBytes:  int64(len(input)),
		OutputBytes: int64(len(out)),
		DurationMs:  float64(elapsed.Microseconds()) / 1000.0,
		Tool:        "base64.decode",
	}, nil
}

func encoder(alphabet Alphabet, padding bool) (*base64.Encoding, error) {
	var enc *base64.Encoding
	switch alphabet {
	case "", AlphabetStandard:
		enc = base64.StdEncoding
	case AlphabetURL:
		enc = base64.URLEncoding
	default:
		return nil, apperr.InvalidRequest("alphabet must be standard or url")
	}
	if !padding {
		enc = enc.WithPadding(base64.NoPadding)
	}
	return enc, nil
}

func ParseAlphabet(s string) (Alphabet, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "", "standard":
		return AlphabetStandard, nil
	case "url":
		return AlphabetURL, nil
	default:
		return "", apperr.InvalidRequest("alphabet must be standard or url")
	}
}
