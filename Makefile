.PHONY: generate generate-discoverability build test lint fmt check validate-spec verify-discoverability run docker-up docker-down tidy install-tools

GO ?= go
OAPI_CODEGEN ?= oapi-codegen
BINARY := toolinfra
VERSION ?= 0.1.0

generate:
	$(OAPI_CODEGEN) -config api/oapi-codegen.yaml api/openapi.yaml

generate-discoverability:
	npm install yaml --no-save 2>/dev/null || true
	node scripts/generate-llms.js
	node -e "const fs=require('fs');const yaml=require('yaml');const d=yaml.parse(fs.readFileSync('api/openapi.yaml','utf8'));fs.writeFileSync('internal/discoverability/spec/openapi.json',JSON.stringify(d,null,2));"
	cp docs/site/llms.txt internal/discoverability/spec/llms.txt
	node -e "const fs=require('fs');const p=['# demo-mcp\n\n',fs.readFileSync('docs/site/llms.txt','utf8'),'\n---\n\n',fs.readFileSync('docs/tools/base64.md','utf8'),'\n---\n\n',fs.readFileSync('docs/errors.md','utf8'),'\n---\n\n',fs.readFileSync('docs/mcp-setup.md','utf8')];fs.writeFileSync('internal/discoverability/spec/llms-full.txt',p.join(''));fs.copyFileSync('internal/discoverability/spec/llms-full.txt','docs/site/llms-full.txt');"

build:
	$(GO) build -ldflags "-X main.version=$(VERSION)" -o bin/$(BINARY) ./cmd/toolinfra

test:
	$(GO) test ./... -count=1

test-unit:
	$(GO) test ./internal/... ./cmd/... ./evals/... -count=1

test-integration:
	$(GO) test ./tests/integration/... -count=1 -v

fmt:
	$(GO) fmt ./...

tidy:
	$(GO) mod tidy

lint:
	golangci-lint run ./...

validate-spec:
	@test -f api/openapi.yaml
	@test -f internal/discoverability/spec/openapi.json

verify-discoverability:
	node scripts/verify_discoverability.mjs

check: validate-spec fmt test-unit test-integration verify-discoverability

run:
	$(GO) run ./cmd/toolinfra serve --api --mcp-http

docker-up:
	docker compose -f deployments/docker-compose.yml up --build -d

docker-down:
	docker compose -f deployments/docker-compose.yml down

install-tools:
	$(GO) install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
