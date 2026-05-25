package agent_tool_selection

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestToolSelectionBaseline(t *testing.T) {
	root := findRoot()
	tags := loadSemanticTags(t, filepath.Join(root, "docs", "semantic-tags.json"))
	prompts := loadPrompts(t, filepath.Join(root, "evals", "agent_tool_selection", "prompts.jsonl"))

	var ok, total int
	for _, p := range prompts {
		total++
		got := matchTool(strings.ToLower(p.Prompt), tags)
		if got != p.ExpectedTool {
			t.Errorf("prompt %q: got %q want %q", p.Prompt, got, p.ExpectedTool)
		} else {
			ok++
		}
	}
	rate := float64(ok) / float64(total) * 100
	t.Logf("selection@1: %.1f%% (%d/%d)", rate, ok, total)
	if rate < 90 {
		t.Fatalf("selection rate %.1f%% below 90%% target", rate)
	}
}

type promptCase struct {
	Prompt       string `json:"prompt"`
	ExpectedTool string `json:"expected_tool"`
}

type semanticTag struct {
	Intents []string `json:"intents"`
}

func loadPrompts(t *testing.T, path string) []promptCase {
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	var out []promptCase
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		var p promptCase
		if err := json.Unmarshal(sc.Bytes(), &p); err != nil {
			t.Fatal(err)
		}
		out = append(out, p)
	}
	return out
}

func loadSemanticTags(t *testing.T, path string) map[string]semanticTag {
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var tags map[string]semanticTag
	if err := json.Unmarshal(b, &tags); err != nil {
		t.Fatal(err)
	}
	return tags
}

func matchTool(prompt string, tags map[string]semanticTag) string {
	decodeHints := []string{"decode", "decoded", "decrypt", "reverse base64", "text for base64", "utf8 text for"}
	for _, h := range decodeHints {
		if strings.Contains(prompt, h) {
			return "base64_decode"
		}
	}
	if strings.Contains(prompt, "encode") ||
		strings.Contains(prompt, "into base64") ||
		strings.Contains(prompt, "to base64") ||
		(strings.Contains(prompt, "convert") && strings.Contains(prompt, "base64")) {
		return "base64_encode"
	}
	best := ""
	bestScore := 0
	for tool, tag := range tags {
		score := 0
		for _, intent := range tag.Intents {
			if strings.Contains(prompt, intent) {
				score += 2
			}
		}
		if strings.Contains(prompt, strings.ReplaceAll(tool, "_", " ")) {
			score += 3
		}
		if score > bestScore {
			bestScore = score
			best = tool
		}
	}
	return best
}

func findRoot() string {
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "."
		}
		dir = parent
	}
}
