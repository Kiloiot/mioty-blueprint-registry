package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRun_FramesProblemsWithStopCommands(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, " ::error::x.json"), []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	if code := run([]string{root}, &stdout, &stderr); code != 1 {
		t.Fatalf("exit code %d, want 1", code)
	}
	lines := strings.Split(strings.TrimSuffix(stdout.String(), "\n"), "\n")
	if len(lines) != 3 {
		t.Fatalf("want stop line, one problem, resume line; got:\n%s", stdout.String())
	}
	token, ok := strings.CutPrefix(lines[0], "::stop-commands::")
	if !ok || len(token) < 16 {
		t.Fatalf("first line must stop workflow commands with a random token, got %q", lines[0])
	}
	if lines[2] != "::"+token+"::" {
		t.Errorf("last line must resume commands with the token, got %q", lines[2])
	}
	if !strings.HasPrefix(lines[1], " ::error::x.json: path: ") {
		t.Errorf("problem line must be inside the frame, got %q", lines[1])
	}
}

func TestRun_ExitCodes(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := run([]string{t.TempDir()}, &stdout, &stderr); code != 0 || stdout.Len() != 0 {
		t.Errorf("clean registry: exit %d, output %q; want 0 and no output", code, stdout.String())
	}
	if code := run(nil, &stdout, &stderr); code != 2 {
		t.Errorf("no arguments: exit %d, want 2", code)
	}
	if code := run([]string{"a", "b"}, &stdout, &stderr); code != 2 {
		t.Errorf("two arguments: exit %d, want 2", code)
	}
}
