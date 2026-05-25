package discoverability

import "testing"

func TestDefaultCatalog(t *testing.T) {
	cat := DefaultCatalog("http://test.example")
	if cat.Version != "0.1.0" {
		t.Fatalf("version: got %q", cat.Version)
	}
	if cat.OpenAPI != "http://test.example/openapi.json" {
		t.Fatalf("openapi url: got %q", cat.OpenAPI)
	}
	if len(cat.Tools) != 2 {
		t.Fatalf("tools: got %d", len(cat.Tools))
	}
	names := map[string]bool{}
	for _, tool := range cat.Tools {
		names[tool.MCPToolName] = true
		if !tool.Deterministic {
			t.Fatalf("%s not deterministic", tool.MCPToolName)
		}
	}
	if !names["base64_encode"] || !names["base64_decode"] {
		t.Fatalf("missing tools: %v", names)
	}
}

func TestDefaultCatalogEmptyBaseURL(t *testing.T) {
	cat := DefaultCatalog("")
	if cat.OpenAPI != "/openapi.json" {
		t.Fatalf("openapi: got %q", cat.OpenAPI)
	}
}
