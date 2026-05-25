package discoverability

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func newDiscTestMux() *chi.Mux {
	r := chi.NewRouter()
	NewHandler("http://test.example").RegisterRoutes(r)
	return r
}

func TestListTools(t *testing.T) {
	w := httptest.NewRecorder()
	newDiscTestMux().ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/v1/tools", nil))

	if w.Code != http.StatusOK {
		t.Fatalf("status: got %d", w.Code)
	}
	var cat ToolCatalog
	if err := json.NewDecoder(w.Body).Decode(&cat); err != nil {
		t.Fatal(err)
	}
	if len(cat.Tools) != 2 {
		t.Fatalf("tools: got %d", len(cat.Tools))
	}
}

func TestServeOpenAPIJSON(t *testing.T) {
	w := httptest.NewRecorder()
	newDiscTestMux().ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/openapi.json", nil))

	if w.Code != http.StatusOK {
		t.Fatalf("status: got %d", w.Code)
	}
	if !strings.Contains(w.Header().Get("Content-Type"), "application/json") {
		t.Fatalf("content-type: %q", w.Header().Get("Content-Type"))
	}
	if len(w.Body.Bytes()) < 10 || w.Body.Bytes()[0] != '{' {
		t.Fatalf("expected JSON object body")
	}
	if !strings.Contains(w.Header().Get("Link"), "llms-txt") {
		t.Fatalf("missing Link header for llms.txt")
	}
}

func TestServeOpenAPIYAML(t *testing.T) {
	w := httptest.NewRecorder()
	newDiscTestMux().ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/openapi.yaml", nil))

	if w.Code != http.StatusOK {
		t.Fatalf("status: got %d", w.Code)
	}
	if !strings.HasPrefix(string(w.Body.Bytes()), "openapi:") {
		t.Fatalf("expected yaml openapi prefix")
	}
}

func TestServeLLMsTxt(t *testing.T) {
	w := httptest.NewRecorder()
	newDiscTestMux().ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/llms.txt", nil))

	if w.Code != http.StatusOK {
		t.Fatalf("status: got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "demo-mcp") {
		t.Fatalf("body missing project name")
	}
}

func TestServeAPICatalog(t *testing.T) {
	w := httptest.NewRecorder()
	newDiscTestMux().ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/.well-known/api-catalog", nil))

	if w.Code != http.StatusOK {
		t.Fatalf("status: got %d", w.Code)
	}
	var doc map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&doc); err != nil {
		t.Fatal(err)
	}
	apis, ok := doc["apis"].([]interface{})
	if !ok || len(apis) == 0 {
		t.Fatal("expected apis array")
	}
}
