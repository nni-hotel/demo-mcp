package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/nni-hotel/demo-mcp/internal/api/gen"
	apperr "github.com/nni-hotel/demo-mcp/internal/platform/errors"
	"github.com/nni-hotel/demo-mcp/internal/platform/middleware"
	"github.com/nni-hotel/demo-mcp/internal/tools/base64"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "toolinfra_http_requests_total",
			Help: "Total HTTP requests",
		},
		[]string{"operation", "status"},
	)
)

func init() {
	prometheus.MustRegister(httpRequests)
}

type Handler struct {
	MaxBodyBytes int64
}

func NewHandler(maxBody int64) *Handler {
	return &Handler{MaxBodyBytes: maxBody}
}

func (h *Handler) Base64Encode(w http.ResponseWriter, r *http.Request) {
	h.handleBase64(w, r, true)
}

func (h *Handler) Base64Decode(w http.ResponseWriter, r *http.Request) {
	h.handleBase64(w, r, false)
}

func (h *Handler) handleBase64(w http.ResponseWriter, r *http.Request, encode bool) {
	op := "base64Decode"
	if encode {
		op = "base64Encode"
	}
	req, err := gen.ReadBase64Request(r)
	if err != nil {
		h.writeError(w, r, apperr.InvalidRequest(err.Error()))
		httpRequests.WithLabelValues(op, "400").Inc()
		return
	}
	toolReq, err := toToolRequest(req)
	if err != nil {
		h.writeError(w, r, err)
		httpRequests.WithLabelValues(op, "400").Inc()
		return
	}
	var result base64.Result
	if encode {
		result, err = base64.Encode(toolReq)
	} else {
		result, err = base64.Decode(toolReq)
	}
	if err != nil {
		h.writeError(w, r, err)
		httpRequests.WithLabelValues(op, statusLabel(err)).Inc()
		return
	}
	resp := gen.ToolResponse{
		Data: gen.Base64Data{Output: result.Output},
		Meta: gen.ToolMeta{
			Tool:        result.Tool,
			InputBytes:  result.InputBytes,
			OutputBytes: result.OutputBytes,
			DurationMs:  result.DurationMs,
		},
	}
	gen.WriteJSON(w, http.StatusOK, resp)
	httpRequests.WithLabelValues(op, "200").Inc()
}

func toToolRequest(req gen.Base64Request) (base64.Request, error) {
	alphabet := "standard"
	if req.Alphabet != nil {
		alphabet = string(*req.Alphabet)
	}
	a, err := base64.ParseAlphabet(alphabet)
	if err != nil {
		return base64.Request{}, err
	}
	padding := true
	if req.Padding != nil {
		padding = *req.Padding
	}
	return base64.Request{
		Input:    req.Input,
		Alphabet: a,
		Padding:  padding,
	}, nil
}

func (h *Handler) Healthz(w http.ResponseWriter, r *http.Request) {
	gen.WriteJSON(w, http.StatusOK, gen.HealthResponse{Status: "ok"})
}

func (h *Handler) Readyz(w http.ResponseWriter, r *http.Request) {
	gen.WriteJSON(w, http.StatusOK, gen.HealthResponse{Status: "ok"})
}

func (h *Handler) Metrics(w http.ResponseWriter, r *http.Request) {
	promhttp.Handler().ServeHTTP(w, r)
}

func (h *Handler) writeError(w http.ResponseWriter, r *http.Request, err error) {
	var ae *apperr.AppError
	if !errors.As(err, &ae) {
		ae = apperr.Internal("unexpected error")
	}
	if errors.Is(err, io.EOF) || isMaxBytesError(err) {
		ae = apperr.PayloadTooLarge(h.MaxBodyBytes)
	}
	rid := middleware.RequestIDFromContext(r.Context())
	resp := gen.ErrorResponse{
		Error: gen.ErrorBody{
			Code:    string(ae.Code),
			Message: ae.Message,
			Details: ae.Details,
		},
		Meta: &gen.ErrorMeta{RequestId: &rid},
	}
	gen.WriteJSON(w, ae.Status, resp)
}

func isMaxBytesError(err error) bool {
	var maxErr *http.MaxBytesError
	return errors.As(err, &maxErr)
}

func statusLabel(err error) string {
	return http.StatusText(apperr.StatusFor(err))
}

// Ensure json import used when extending handlers.
var _ = json.Marshal
