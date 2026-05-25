package integration

import (
	"net/http"
	"strings"
	"testing"

	"github.com/nni-hotel/demo-mcp/internal/api/gen"
	"github.com/nni-hotel/demo-mcp/internal/discoverability"
	"github.com/nni-hotel/demo-mcp/internal/testutil"
)

func TestIntegration_Base64EncodeDecodeRoundTrip(t *testing.T) {
	srv := testutil.NewTestServer(t)
	defer srv.Close()

	encResp := testutil.DoJSON(t, srv, http.MethodPost, "/v1/tools/base64/encode", map[string]interface{}{
		"input":    "integration test",
		"alphabet": "standard",
		"padding":  true,
	})
	if encResp.StatusCode != http.StatusOK {
		t.Fatalf("encode status: %d", encResp.StatusCode)
	}
	var enc gen.ToolResponse
	testutil.DecodeJSON(t, encResp, &enc)

	decResp := testutil.DoJSON(t, srv, http.MethodPost, "/v1/tools/base64/decode", map[string]interface{}{
		"input":    enc.Data.Output,
		"alphabet": "standard",
		"padding":  true,
	})
	if decResp.StatusCode != http.StatusOK {
		t.Fatalf("decode status: %d", decResp.StatusCode)
	}
	var dec gen.ToolResponse
	testutil.DecodeJSON(t, decResp, &dec)

	if dec.Data.Output != "integration test" {
		t.Fatalf("round-trip: got %q", dec.Data.Output)
	}
}

func TestIntegration_DiscoverabilityEndpoints(t *testing.T) {
	srv := testutil.NewTestServer(t)
	defer srv.Close()

	tests := []struct {
		path       string
		status     int
		bodySubstr string
	}{
		{"/healthz", http.StatusOK, `"ok"`},
		{"/readyz", http.StatusOK, `"ok"`},
		{"/v1/tools", http.StatusOK, "base64_encode"},
		{"/openapi.json", http.StatusOK, "openapi"},
		{"/openapi.yaml", http.StatusOK, "openapi:"},
		{"/llms.txt", http.StatusOK, "demo-mcp"},
		{"/.well-known/api-catalog", http.StatusOK, "demo-mcp"},
	}
	for _, tc := range tests {
		t.Run(tc.path, func(t *testing.T) {
			resp, err := http.Get(srv.URL + tc.path)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != tc.status {
				t.Fatalf("status: got %d", resp.StatusCode)
			}
			buf := make([]byte, 4096)
			n, _ := resp.Body.Read(buf)
			body := string(buf[:n])
			if tc.bodySubstr != "" && !strings.Contains(body, tc.bodySubstr) {
				t.Fatalf("body missing %q: %s", tc.bodySubstr, truncate(body, 200))
			}
		})
	}
}

func TestIntegration_ToolCatalogMatchesDiscoverability(t *testing.T) {
	srv := testutil.NewTestServer(t)
	defer srv.Close()

	resp := testutil.DoJSON(t, srv, http.MethodGet, "/v1/tools", nil)
	defer resp.Body.Close()

	var cat discoverability.ToolCatalog
	testutil.DecodeJSON(t, resp, &cat)

	expected := discoverability.DefaultCatalog("http://test.example")
	if len(cat.Tools) != len(expected.Tools) {
		t.Fatalf("tool count: got %d want %d", len(cat.Tools), len(expected.Tools))
	}
	for i, tool := range cat.Tools {
		if tool.MCPToolName != expected.Tools[i].MCPToolName {
			t.Fatalf("tool[%d] name: got %q want %q", i, tool.MCPToolName, expected.Tools[i].MCPToolName)
		}
	}
}

func TestIntegration_RequestIDHeader(t *testing.T) {
	srv := testutil.NewTestServer(t)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/healthz")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.Header.Get("X-Request-ID") == "" {
		t.Fatal("missing X-Request-ID header")
	}
}

func TestIntegration_PayloadTooLarge(t *testing.T) {
	srv := testutil.NewTestServer(t)
	defer srv.Close()

	huge := map[string]string{"input": strings.Repeat("x", 2_000_000)}
	resp := testutil.DoJSON(t, srv, http.MethodPost, "/v1/tools/base64/encode", huge)
	if resp.StatusCode != http.StatusRequestEntityTooLarge {
		t.Fatalf("status: got %d", resp.StatusCode)
	}
}

func TestIntegration_InvalidBase64ReturnsStructuredError(t *testing.T) {
	srv := testutil.NewTestServer(t)
	defer srv.Close()

	resp := testutil.DoJSON(t, srv, http.MethodPost, "/v1/tools/base64/decode", map[string]string{
		"input": "@@@",
	})
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("status: got %d", resp.StatusCode)
	}
	var errResp gen.ErrorResponse
	testutil.DecodeJSON(t, resp, &errResp)
	if errResp.Error.Code != "INVALID_BASE64" {
		t.Fatalf("code: got %q", errResp.Error.Code)
	}
	if errResp.Meta == nil || errResp.Meta.RequestId == nil || *errResp.Meta.RequestId == "" {
		t.Fatal("expected request_id in error meta")
	}
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
