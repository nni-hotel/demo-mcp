#!/usr/bin/env node
import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const root = path.join(path.dirname(fileURLToPath(import.meta.url)), '..');
const failures = [];

function read(p) {
  return fs.readFileSync(path.join(root, p), 'utf8');
}

const openapi = read('api/openapi.yaml');
const catalog = read('internal/discoverability/catalog.go');
const tools = JSON.parse(read('docs/mcp/tools.json'));

for (const t of tools.tools) {
  if (!openapi.includes(`operationId: ${t.openapi_operation_id}`)) {
    failures.push(`openapi missing operationId: ${t.openapi_operation_id}`);
  }
  if (!openapi.includes(`x-mcp-tool-name: ${t.name}`)) {
    failures.push(`openapi missing x-mcp-tool-name: ${t.name}`);
  }
  if (!catalog.includes(t.name)) {
    failures.push(`catalog.go missing MCP tool: ${t.name}`);
  }
}

for (const f of [
  'internal/discoverability/spec/openapi.json',
  'internal/discoverability/spec/llms.txt',
  'internal/discoverability/spec/llms-full.txt',
  'server.json',
  'docs/semantic-tags.yaml',
]) {
  if (!fs.existsSync(path.join(root, f))) failures.push(`missing ${f}`);
}

const tagsYaml = read('docs/semantic-tags.yaml');
for (const t of tools.tools) {
  if (!tagsYaml.includes(`${t.name}:`)) {
    failures.push(`semantic-tags.yaml missing key: ${t.name}`);
  }
}

if (failures.length) {
  console.error('verify-discoverability failed:');
  failures.forEach((f) => console.error('  -', f));
  process.exit(1);
}
console.log('verify-discoverability: OK');
