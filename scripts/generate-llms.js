const fs = require('fs');
const path = require('path');

const root = path.join(__dirname, '..');
const site = path.join(root, 'docs', 'site');
const spec = path.join(root, 'internal', 'discoverability', 'spec');

const llms = fs.readFileSync(path.join(site, 'llms.txt'), 'utf8');
const parts = [
  '# demo-mcp (full)\n\n',
  llms,
  '\n---\n\n',
  fs.readFileSync(path.join(root, 'docs', 'tools', 'base64.md'), 'utf8'),
  '\n---\n\n',
  fs.readFileSync(path.join(root, 'docs', 'errors.md'), 'utf8'),
  '\n---\n\n',
  fs.readFileSync(path.join(root, 'docs', 'mcp-setup.md'), 'utf8'),
];
const full = parts.join('');
fs.writeFileSync(path.join(site, 'llms-full.txt'), full);
fs.mkdirSync(spec, { recursive: true });
fs.copyFileSync(path.join(site, 'llms.txt'), path.join(spec, 'llms.txt'));
fs.writeFileSync(path.join(spec, 'llms-full.txt'), full);
console.log('generated llms.txt and llms-full.txt');
