package mcp

import (
	"context"
	"encoding/json"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/nni-hotel/demo-mcp/internal/config"
	"github.com/nni-hotel/demo-mcp/internal/tools/base64"
)

type ToolInput struct {
	Input    string `json:"input" jsonschema:"UTF-8 text to encode or Base64 string to decode"`
	Alphabet string `json:"alphabet,omitempty" jsonschema:"standard or url,default standard"`
	Padding  *bool  `json:"padding,omitempty" jsonschema:"include padding,default true"`
}

type ToolOutput struct {
	Data ToolOutputData `json:"data"`
	Meta ToolOutputMeta `json:"meta"`
}

type ToolOutputData struct {
	Output string `json:"output"`
}

type ToolOutputMeta struct {
	Tool        string  `json:"tool"`
	InputBytes  int64   `json:"input_bytes"`
	OutputBytes int64   `json:"output_bytes"`
	DurationMs  float64 `json:"duration_ms"`
}

func NewServer(cfg config.Config) *mcp.Server {
	s := mcp.NewServer(&mcp.Implementation{
		Name:    cfg.ServiceName,
		Version: cfg.Version,
	}, nil)

	mcp.AddTool(s, &mcp.Tool{
		Name:        "base64_encode",
		Description: "Encode UTF-8 text to Base64 (standard or URL-safe). Deterministic. Not for binary.",
	}, encodeHandler)
	mcp.AddTool(s, &mcp.Tool{
		Name:        "base64_decode",
		Description: "Decode Base64 to UTF-8 text. Returns INVALID input on bad Base64. Not for binary.",
	}, decodeHandler)

	return s
}

func encodeHandler(ctx context.Context, req *mcp.CallToolRequest, in ToolInput) (*mcp.CallToolResult, ToolOutput, error) {
	return runTool(in, true)
}

func decodeHandler(ctx context.Context, req *mcp.CallToolRequest, in ToolInput) (*mcp.CallToolResult, ToolOutput, error) {
	return runTool(in, false)
}

func runTool(in ToolInput, encode bool) (*mcp.CallToolResult, ToolOutput, error) {
	a, err := base64.ParseAlphabet(in.Alphabet)
	if err != nil {
		return toolError(err)
	}
	padding := true
	if in.Padding != nil {
		padding = *in.Padding
	}
	treq := base64.Request{Input: in.Input, Alphabet: a, Padding: padding}
	var res base64.Result
	if encode {
		res, err = base64.Encode(treq)
	} else {
		res, err = base64.Decode(treq)
	}
	if err != nil {
		return toolError(err)
	}
	out := ToolOutput{
		Data: ToolOutputData{Output: res.Output},
		Meta: ToolOutputMeta{
			Tool:        res.Tool,
			InputBytes:  res.InputBytes,
			OutputBytes: res.OutputBytes,
			DurationMs:  res.DurationMs,
		},
	}
	return nil, out, nil
}

func toolError(err error) (*mcp.CallToolResult, ToolOutput, error) {
	b, _ := json.Marshal(map[string]string{"error": err.Error()})
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: string(b)}},
		IsError: true,
	}, ToolOutput{}, nil
}
