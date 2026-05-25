package discoverability

import _ "embed"

//go:embed spec/openapi.yaml
var openAPIYAML []byte

//go:embed spec/openapi.json
var openAPIJSON []byte

//go:embed spec/llms.txt
var defaultLLMsTxt []byte

//go:embed spec/llms-full.txt
var defaultLLMsFullTxt []byte

// EmbeddedOpenAPIYAML returns the OpenAPI specification bytes.
func EmbeddedOpenAPIYAML() []byte {
	return openAPIYAML
}

// DefaultLLMsTxt returns embedded llms.txt (sync via make generate-discoverability).
func DefaultLLMsTxt() []byte {
	return defaultLLMsTxt
}

// DefaultLLMsFullTxt returns embedded llms-full.txt.
func DefaultLLMsFullTxt() []byte {
	return defaultLLMsFullTxt
}
