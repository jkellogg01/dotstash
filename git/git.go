package git

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/log"
)

func Download(url, branch, alias string) error {
	args := []string{
		"clone",
		"--depth=1",
	}
	if branch != "" {
		args = append(args, "--branch="+branch)
	}
	args = append(args, url)
	if alias != "" {
		args = append(args, alias)
	}
	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	log.Debug("", "output", output)
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
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	err = os.Chdir(path)
	if err != nil {
		return err
	}
	cmd := exec.Command("git", "init")
	log.Debugf("running command %s", cmd.String())
	buf, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	log.Debug("", "output", buf)
	return os.Chdir(wd)
}
