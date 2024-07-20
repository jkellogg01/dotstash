package files

import (
	"os"
	"strings"
	"testing"
)

func TestGetFigurePath(t *testing.T) {
	path, err := GetFigurePath()
	if err != nil {
		t.Fatal(err)
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get user home dir")
	}
	_, cut := strings.CutPrefix(path, homeDir)
	if !cut {
		t.Error("path does not contain home directory")
	}
	segments := strings.Split(path, string(os.PathSeparator))
	if len(segments) <= 0 {
		t.Errorf("got malformed path: %s", path)
	} else if segments[len(segments)-1] != ".figure" {
		t.Error("path does not lead to a '.figure' folder")
	}
}
