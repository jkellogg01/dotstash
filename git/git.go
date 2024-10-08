package git

import (
	"io"
	"os/exec"
	"strings"

	"github.com/charmbracelet/log"
)

func Download(url, branch, dst string) error {
	args := []string{"clone"}
	if branch != "" {
		args = append(args, "--branch="+branch)
	}
	args = append(args, url)
	if dst != "" {
		args = append(args, dst)
	}
	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	log.Debug("", "output", string(output))
	return err
}

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

func InitRepo(path string) error {
	cmd := exec.Command("git", "init")
	cmd.Dir = path
	log.Debugf("running command %s", cmd.String())
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	log.Debug("", "output", string(output))
	return nil
}
