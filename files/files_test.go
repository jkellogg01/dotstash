package files

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var homePath string

func TestGetDotstashPath(t *testing.T) {
	path, err := GetDotstashPath()
	if err != nil {
		t.Fatal(err)
	}
	_, cut := strings.CutPrefix(path, homePath)
	if !cut {
		t.Error("path does not contain home directory")
	}
	segments := filepath.SplitList(path)
	if len(segments) == 0 {
		t.Errorf("got malformed path: %s", path)
	}
	if filepath.Base(path) != ".dotstash" {
		t.Error("path does not lead to a '.dotstash' folder")
	}
}

func init() {
	var err error
	homePath, err = os.UserHomeDir()
	if err != nil {
		panic(err)
	}
}
