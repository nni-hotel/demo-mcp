package mcp

import (
	"context"
	"testing"

	"github.com/nni-hotel/demo-mcp/internal/config"
)

func testCfg() config.Config {
	return config.Config{
		ServiceName: "demo-mcp-test",
		Version:     "0.1.0-test",
	}
}

func TestRunToolEncode(t *testing.T) {
	_, out, err := runTool(ToolInput{Input: "hello"}, true)
	if err != nil {
		t.Fatal(err)
	}
	if out.Data.Output != "aGVsbG8=" {
		t.Fatalf("output: got %q", out.Data.Output)
	}
	if out.Meta.Tool != "base64.encode" {
		t.Fatalf("tool: got %q", out.Meta.Tool)
	}
}

func TestRunToolDecode(t *testing.T) {
	_, out, err := runTool(ToolInput{Input: "aGVsbG8="}, false)
	if err != nil {
		t.Fatal(err)
	}
	if out.Data.Output != "hello" {
		t.Fatalf("output: got %q", out.Data.Output)
	}
}

func TestRunToolInvalidBase64(t *testing.T) {
	result, _, err := runTool(ToolInput{Input: "not!!!"}, false)
	if err != nil {
		t.Fatal(err)
	}
	if result == nil || !result.IsError {
		t.Fatal("expected MCP error result")
	}
}

func TestEncodeHandler(t *testing.T) {
	_, out, err := encodeHandler(context.Background(), nil, ToolInput{Input: "hi"})
	if err != nil {
		t.Fatal(err)
	}
	if out.Data.Output == "" {
		t.Fatal("expected non-empty output")
	}
}

func TestNewServerRegistersTools(t *testing.T) {
	s := NewServer(testCfg())
	if s == nil {
		t.Fatal("nil server")
	}
}
