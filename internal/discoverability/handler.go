package discoverability

import (
	"net/http"
	"strings"

	"github.com/nni-hotel/demo-mcp/internal/api/gen"
)

// Handler serves discoverability endpoints (OpenAPI, catalog, llms.txt, well-known).
type Handler struct {
	OpenAPIYAML []byte
	OpenAPIJSON []byte
	LLMsTxt     []byte
	LLMsFullTxt []byte
	BaseURL     string
}

func NewHandler(baseURL string) *Handler {
	return &Handler{
		OpenAPIYAML: EmbeddedOpenAPIYAML(),
		OpenAPIJSON: OpenAPIJSONBytes(),
		LLMsTxt:     DefaultLLMsTxt(),
		LLMsFullTxt: DefaultLLMsFullTxt(),
		BaseURL:     strings.TrimRight(baseURL, "/"),
	}
}

func (h *Handler) RegisterRoutes(mux interface {
	Get(pattern string, handlerFn http.HandlerFunc)
}) {
	mux.Get("/v1/tools", h.ListTools)
	mux.Get("/openapi.json", h.ServeOpenAPIJSON)
	mux.Get("/openapi.yaml", h.ServeOpenAPIYAML)
	mux.Get("/llms.txt", h.ServeLLMsTxt)
	mux.Get("/llms-full.txt", h.ServeLLMsFullTxt)
	mux.Get("/.well-known/api-catalog", h.ServeAPICatalog)
}

func (h *Handler) ListTools(w http.ResponseWriter, r *http.Request) {
	gen.WriteJSON(w, http.StatusOK, DefaultCatalog(h.baseURL(r)))
}

func (h *Handler) ServeOpenAPIJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Link", `<`+h.baseURL(r)+`/llms.txt>; rel="llms-txt"`)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(h.OpenAPIJSON)
}

func (h *Handler) ServeOpenAPIYAML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/yaml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(h.OpenAPIYAML)
}

func (h *Handler) ServeLLMsTxt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(h.LLMsTxt)
}

func (h *Handler) ServeLLMsFullTxt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(h.LLMsFullTxt)
}

func (h *Handler) ServeAPICatalog(w http.ResponseWriter, r *http.Request) {
	base := h.baseURL(r)
	catalog := map[string]interface{}{
		"version": "0.1.0",
		"apis": []map[string]interface{}{
			{
				"name":        "demo-mcp",
				"description": "Deterministic AI utility tools (Base64 v0.1)",
				"openapi":     base + "/openapi.json",
				"catalog":     base + "/v1/tools",
				"mcp": map[string]string{
					"registry": "io.github.nni-hotel/demo-mcp",
					"stdio":    "toolinfra mcp stdio",
				},
			},
		},
	}
	gen.WriteJSON(w, http.StatusOK, catalog)
}

func (h *Handler) baseURL(r *http.Request) string {
	if h.BaseURL != "" {
		return h.BaseURL
	}
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + r.Host
}
