//go:build ignore

// Generates internal/discoverability/spec/openapi.json from api/openapi.yaml
// and copies llms.txt from docs/site/.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func main() {
	root := filepath.Join("..")
	if _, err := os.Stat("api/openapi.yaml"); err == nil {
		root = "."
	}
	yamlPath := filepath.Join(root, "api", "openapi.yaml")
	jsonPath := filepath.Join(root, "internal", "discoverability", "spec", "openapi.json")
	llmsSrc := filepath.Join(root, "docs", "site", "llms.txt")
	llmsDst := filepath.Join(root, "internal", "discoverability", "spec", "llms.txt")
	llmsFullSrc := filepath.Join(root, "docs", "site", "llms-full.txt")
	llmsFullDst := filepath.Join(root, "internal", "discoverability", "spec", "llms-full.txt")

	yamlBytes, err := os.ReadFile(yamlPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read openapi yaml: %v\n", err)
		os.Exit(1)
	}
	var doc interface{}
	if err := yaml.Unmarshal(yamlBytes, &doc); err != nil {
		fmt.Fprintf(os.Stderr, "parse yaml: %v\n", err)
		os.Exit(1)
	}
	jsonBytes, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "marshal json: %v\n", err)
		os.Exit(1)
	}
	if err := os.MkdirAll(filepath.Dir(jsonPath), 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "mkdir: %v\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile(jsonPath, jsonBytes, 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "write json: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("wrote", jsonPath)

	copyFile(llmsSrc, llmsDst)
	copyFile(llmsFullSrc, llmsFullDst)
}

func copyFile(src, dst string) {
	b, err := os.ReadFile(src)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read %s: %v\n", src, err)
		os.Exit(1)
	}
	if err := os.WriteFile(dst, b, 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "write %s: %v\n", dst, err)
		os.Exit(1)
	}
	fmt.Println("wrote", dst)
}
