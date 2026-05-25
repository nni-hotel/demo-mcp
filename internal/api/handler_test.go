package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nni-hotel/demo-mcp/internal/api/gen"
)

func TestBase64EncodeHandler(t *testing.T) {
	h := NewHandler(1 << 20)
	body := `{"input":"hello"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/tools/base64/encode", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Base64Encode(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status: got %d body %s", w.Code, w.Body.String())
	}
	var resp gen.ToolResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.Data.Output != "aGVsbG8=" {
		t.Fatalf("output: got %q", resp.Data.Output)
	}
	if resp.Meta.Tool != "base64.encode" {
		t.Fatalf("meta.tool: got %q", resp.Meta.Tool)
	}
}

func TestBase64DecodeHandler(t *testing.T) {
	h := NewHandler(1 << 20)
	body := `{"input":"aGVsbG8="}`
	req := httptest.NewRequest(http.MethodPost, "/v1/tools/base64/decode", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Base64Decode(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status: got %d", w.Code)
	}
	var resp gen.ToolResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.Data.Output != "hello" {
		t.Fatalf("output: got %q", resp.Data.Output)
	}
}

func TestBase64DecodeInvalid(t *testing.T) {
	h := NewHandler(1 << 20)
	req := httptest.NewRequest(http.MethodPost, "/v1/tools/base64/decode", strings.NewReader(`{"input":"!!!"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Base64Decode(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d", w.Code)
	}
	var resp gen.ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.Error.Code != "INVALID_BASE64" {
		t.Fatalf("code: got %q", resp.Error.Code)
	}
}

func TestBase64InvalidJSON(t *testing.T) {
	h := NewHandler(1 << 20)
	req := httptest.NewRequest(http.MethodPost, "/v1/tools/base64/encode", strings.NewReader(`{`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Base64Encode(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d", w.Code)
	}
}

func TestBase64InvalidAlphabet(t *testing.T) {
	h := NewHandler(1 << 20)
	req := httptest.NewRequest(http.MethodPost, "/v1/tools/base64/encode",
		strings.NewReader(`{"input":"x","alphabet":"rot13"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Base64Encode(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d", w.Code)
	}
}

func TestPayloadTooLarge(t *testing.T) {
	h := NewHandler(32)
	large := bytes.Repeat([]byte("a"), 64)
	req := httptest.NewRequest(http.MethodPost, "/v1/tools/base64/encode", bytes.NewReader(large))
	req.Header.Set("Content-Type", "application/json")
	req.Body = http.MaxBytesReader(httptest.NewRecorder(), req.Body, 32)
	w := httptest.NewRecorder()

	h.Base64Encode(w, req)

	if w.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status: got %d body %s", w.Code, w.Body.String())
	}
}

func TestHealthz(t *testing.T) {
	h := NewHandler(1 << 20)
	w := httptest.NewRecorder()
	h.Healthz(w, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("status: got %d", w.Code)
	}
}
