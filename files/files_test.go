package files

import (
	"os"
	"testing"
)

func TestGotoFigureRoot(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("could not find user home dir: %s", err)
	}
	err = os.Chdir(home)
	if err != nil {
		t.Fatalf("could not move to home dir: %s", err)
	}
	err = GotoFigureRoot()
	if err != nil {
		t.Fatalf("GotoFigureRoot returned an error: %s", err)
	}
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("could not get working directory: %s", err)
	}
	if wd != home+string(os.PathSeparator)+".figure" {
		t.Fatalf("GotoFigureRoot did not result in wd being ~/.figure")
	}
}
