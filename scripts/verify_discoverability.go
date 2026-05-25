//go:build ignore

// verify_discoverability checks OpenAPI, MCP tools.json, semantic-tags, and catalog alignment.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func main() {
	root := "."
	if _, err := os.Stat("api/openapi.yaml"); os.IsNotExist(err) {
		root = ".."
	}
	var failures []string

	// Load semantic tags
	tagsPath := root + "/docs/semantic-tags.yaml"
	tagsBytes, err := os.ReadFile(tagsPath)
	if err != nil {
		failures = append(failures, "read semantic-tags: "+err.Error())
	} else {
		var tags map[string]map[string]interface{}
		if err := yaml.Unmarshal(tagsBytes, &tags); err != nil {
			failures = append(failures, "parse semantic-tags: "+err.Error())
		} else {
			for name, tag := range tags {
				mcpName, _ := tag["mcp_tool_name"].(string)
				if mcpName != name && mcpName != "" {
					failures = append(failures, fmt.Sprintf("semantic-tags: key %s mcp_tool_name=%s", name, mcpName))
				}
			}
		}
	}

	// Load docs/mcp/tools.json
	toolsPath := root + "/docs/mcp/tools.json"
	toolsBytes, err := os.ReadFile(toolsPath)
	if err != nil {
		failures = append(failures, "read tools.json: "+err.Error())
	} else {
		var doc struct {
			Tools []struct {
				Name                 string   `json:"name"`
				OpenAPIOperationID  string   `json:"openapi_operation_id"`
				SemanticTags         []string `json:"semantic_tags"`
			} `json:"tools"`
		}
		if err := json.Unmarshal(toolsBytes, &doc); err != nil {
			failures = append(failures, "parse tools.json: "+err.Error())
		} else {
			openapiBytes, _ := os.ReadFile(root + "/api/openapi.yaml")
			openapiStr := string(openapiBytes)
			for _, t := range doc.Tools {
				if !strings.Contains(openapiStr, "operationId: "+t.OpenAPIOperationID) {
					failures = append(failures, "openapi missing operationId: "+t.OpenAPIOperationID)
				}
				if !strings.Contains(openapiStr, "x-mcp-tool-name: "+t.Name) {
					failures = append(failures, "openapi missing x-mcp-tool-name: "+t.Name)
				}
			}
		}
	}

	// Catalog MCP names
	catalogPath := root + "/internal/discoverability/catalog.go"
	catalogBytes, _ := os.ReadFile(catalogPath)
	catalogStr := string(catalogBytes)
	for _, mcp := range []string{"base64_encode", "base64_decode"} {
		if !strings.Contains(catalogStr, mcp) {
			failures = append(failures, "catalog.go missing MCP tool: "+mcp)
		}
	}

	// Embedded spec files
	for _, f := range []string{
		"internal/discoverability/spec/openapi.json",
		"internal/discoverability/spec/llms.txt",
		"internal/discoverability/spec/llms-full.txt",
	} {
		if _, err := os.Stat(root + "/" + f); err != nil {
			failures = append(failures, "missing "+f)
		}
	}

	if len(failures) > 0 {
		fmt.Fprintf(os.Stderr, "verify-discoverability failed:\n")
		for _, f := range failures {
			fmt.Fprintf(os.Stderr, "  - %s\n", f)
		}
		os.Exit(1)
	}
	fmt.Println("verify-discoverability: OK")
}
