package discoverability

// ToolCatalog is the machine-readable tool index for agents.
type ToolCatalog struct {
	Version string            `json:"version"`
	OpenAPI string            `json:"openapi"`
	Tools   []ToolCatalogEntry `json:"tools"`
}

type ToolCatalogEntry struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	MCPToolName   string   `json:"mcp_tool_name"`
	RESTPath      string   `json:"rest_path"`
	Method        string   `json:"method"`
	Deterministic bool     `json:"deterministic"`
	SemanticTags  []string `json:"semantic_tags,omitempty"`
	Description   string   `json:"description,omitempty"`
}

// DefaultCatalog returns the v0.1 tool catalog (source of truth for /v1/tools).
func DefaultCatalog(baseURL string) ToolCatalog {
	openapi := baseURL + "/openapi.json"
	if baseURL == "" {
		openapi = "/openapi.json"
	}
	return ToolCatalog{
		Version: "0.1.0",
		OpenAPI: openapi,
		Tools: []ToolCatalogEntry{
			{
				ID:            "base64.encode",
				Name:          "Base64 Encode",
				MCPToolName:   "base64_encode",
				RESTPath:      "/v1/tools/base64/encode",
				Method:        "POST",
				Deterministic: true,
				SemanticTags:  []string{"encoding", "base64", "text", "utf8"},
				Description:   "Encode UTF-8 text to Base64 (standard or URL-safe).",
			},
			{
				ID:            "base64.decode",
				Name:          "Base64 Decode",
				MCPToolName:   "base64_decode",
				RESTPath:      "/v1/tools/base64/decode",
				Method:        "POST",
				Deterministic: true,
				SemanticTags:  []string{"decoding", "base64", "text", "utf8"},
				Description:   "Decode Base64 to UTF-8 text.",
			},
		},
	}
}
