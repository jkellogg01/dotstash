package git

import (
	"os"
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

	err = os.Chdir(dir)
	if err != nil {
		t.Fatal(err)
	}

	c := exec.Command(
		"git",
		"status",
	)
	output, err := c.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains("fatal: not a git repository (or any of the parent directories)", string(output)) {
		t.Fatal("target repository does not contain a git repository")
	}
}
