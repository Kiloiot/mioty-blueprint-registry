package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestWorkflow_PinsTheBuildEnvironment(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "workflows", "check.yml"))
	if err != nil {
		t.Fatal(err)
	}
	wf := string(data)
	env := jobEnv(t, wf)
	for _, want := range []string{"GOFLAGS: -mod=readonly", `GOWORK: "off"`, "GOTOOLCHAIN: local"} {
		if !strings.Contains(env, "\n      "+want+"\n") {
			t.Errorf("job env must contain %q, got:\n%s", want, env)
		}
	}
	for _, want := range []string{"\npermissions:\n  contents: read\n", "\n  pull_request:\n", "\n  check:\n"} {
		if !strings.Contains(wf, want) {
			t.Errorf("workflow must contain %q", want)
		}
	}
	if strings.Contains(wf, "secrets.") {
		t.Error("workflow must not use secrets")
	}
	if strings.Contains(wf, "pull_request_target") {
		t.Error("workflow must not run on pull_request_target")
	}
	for _, uses := range regexp.MustCompile(`uses: (\S+)`).FindAllStringSubmatch(wf, -1) {
		if !regexp.MustCompile(`@[0-9a-f]{40}$`).MatchString(uses[1]) {
			t.Errorf("action %s must be pinned by a full commit SHA", uses[1])
		}
	}
	registry := strings.Index(wf, "run: go run . ../..")
	tests := strings.Index(wf, "run: go test ./...")
	if registry < 0 || tests < 0 || registry > tests {
		t.Error("the registry check must run before go test")
	}
	for _, step := range []string{"go test ./...", "go mod tidy -diff"} {
		if !strings.Contains(wf, "if: ${{ !cancelled() }}\n        working-directory: .github/check\n        run: "+step) {
			t.Errorf("%s must run after a red registry check", step)
		}
	}
}

func jobEnv(t *testing.T, wf string) string {
	t.Helper()
	start := strings.Index(wf, "\n    env:\n")
	if start < 0 {
		t.Fatal("job env block not found")
	}
	block := wf[start+len("\n    env:"):]
	if end := strings.Index(block, "\n    steps:"); end >= 0 {
		return block[:end+1]
	}
	t.Fatal("steps after job env not found")
	return ""
}
