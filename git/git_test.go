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

func TestExecInteractive(t *testing.T) {
	// NOTE: `go test` should just cache the results of TestInitRepo so I don't
	// feel too bad about the way I'm using it here.
	if !t.Run("repo initialization works", TestInitRepo) {
		t.Fatal("this test cannot be run unless TestInitRepo is passing")
	}

	dir := t.TempDir()
	err := InitRepo(dir)
	// NOTE: in this context there should literally never be an error
	// but handling your errors never hurt anybody
	if err != nil {
		t.Fatal(err)
	}

	go func(t *testing.T) {
		err = Exec(dir, "commit")
		if err != nil {
			// I think this should be fatal but testing doesn't like that
			t.Error(err)
		}
	}(t)

	// TODO: script the input so that we can actually close the editor
	// TODO: fail if the editor doesn't close when the script is done
}
