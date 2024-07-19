package git

import (
	"io"
	"os/exec"
	"strings"
)

func CheckGitInstalled() (string, error) {
	cmd := exec.Command("git", "--version")
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	if err := cmd.Start(); err != nil {
		return "", err
	}
	data, err := io.ReadAll(pipe)
	if err != nil {
		return "", err
	}
	if err := cmd.Wait(); err != nil {
		return "", err
	}
	after, _ := strings.CutPrefix(string(data), "git version ")
	return after, nil
}
