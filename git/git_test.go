package git

import (
	"os/exec"
	"strings"
	"testing"
)

func TestInitRepo(t *testing.T) {
	dir := t.TempDir()
	err := InitRepo(dir)
	if err != nil {
		t.Fatalf("repository initiation failed: %v", err)
	}

	c := exec.Command(
		"git",
		"status",
	)
	c.Dir = dir
	output, err := c.CombinedOutput()
	if err != nil {
		t.Error(err)
	}
	if strings.Contains("fatal: not a git repository (or any of the parent directories)", string(output)) {
		t.Error("target repository does not contain a git repository")
	}
}
