// Package gen provides OpenAPI-generated types and Chi server bindings.
// Regenerate with: make generate
package gen

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/oapi-codegen/runtime"
)

const (
	AlphabetStandard Alphabet = "standard"
	AlphabetUrl      Alphabet = "url"
)

type Alphabet string

type Base64Data struct {
	Output string `json:"output"`
}

type Base64Request struct {
	Input    string   `json:"input"`
	Alphabet *Alphabet `json:"alphabet,omitempty"`
	Padding  *bool    `json:"padding,omitempty"`
}

type ErrorBody struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

type ErrorMeta struct {
	RequestId *string `json:"request_id,omitempty"`
}

type ErrorResponse struct {
	Error ErrorBody  `json:"error"`
	Meta  *ErrorMeta `json:"meta,omitempty"`
}

type HealthResponse struct {
	Status string `json:"status"`
}

type ToolMeta struct {
	Tool        string  `json:"tool"`
	InputBytes  int64   `json:"input_bytes"`
	OutputBytes int64   `json:"output_bytes"`
	DurationMs  float64 `json:"duration_ms"`
}

type ToolResponse struct {
	Data Base64Data `json:"data"`
	Meta ToolMeta   `json:"meta"`
}

type ServerInterface interface {
	Base64Encode(w http.ResponseWriter, r *http.Request)
	Base64Decode(w http.ResponseWriter, r *http.Request)
	Healthz(w http.ResponseWriter, r *http.Request)
	Readyz(w http.ResponseWriter, r *http.Request)
	Metrics(w http.ResponseWriter, r *http.Request)
}

type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

func (siw *ServerInterfaceWrapper) Base64Encode(w http.ResponseWriter, r *http.Request) {
	siw.Handler.Base64Encode(w, r)
}

func (siw *ServerInterfaceWrapper) Base64Decode(w http.ResponseWriter, r *http.Request) {
	siw.Handler.Base64Decode(w, r)
}

func (siw *ServerInterfaceWrapper) Healthz(w http.ResponseWriter, r *http.Request) {
	siw.Handler.Healthz(w, r)
}

func (siw *ServerInterfaceWrapper) Readyz(w http.ResponseWriter, r *http.Request) {
	siw.Handler.Readyz(w, r)
}

func (siw *ServerInterfaceWrapper) Metrics(w http.ResponseWriter, r *http.Request) {
	siw.Handler.Metrics(w, r)
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{BaseRouter: r})
}

func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter
	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/v1/tools/base64/encode", wrapper.Base64Encode)
		r.Post(options.BaseURL+"/v1/tools/base64/decode", wrapper.Base64Decode)
		r.Get(options.BaseURL+"/healthz", wrapper.Healthz)
		r.Get(options.BaseURL+"/readyz", wrapper.Readyz)
		r.Get(options.BaseURL+"/metrics", wrapper.Metrics)
	})
	return r
}

// Strict decode helpers used by handlers.

func ReadBase64Request(r *http.Request) (Base64Request, error) {
	var req Base64Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return Base64Request{}, fmt.Errorf("invalid json: %w", err)
	}
	return req, nil
}

func WriteJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// RequestIdFromContext returns request ID if set by middleware.
type requestIDKey struct{}

func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, id)
}

func RequestIDFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(requestIDKey{}).(string); ok {
		return v
	}
	return ""
}

// Unused import guard for runtime package when extending generated code.
var _ = runtime.BindStyledParameter
