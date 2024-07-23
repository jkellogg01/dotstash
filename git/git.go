package git

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/log"
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
	outbuf, errbuf := new(bytes.Buffer), new(bytes.Buffer)
	cmd.Stdout = outbuf
	cmd.Stderr = errbuf
	log.Debugf("running command %s", cmd.String())
	err = cmd.Run()
	if err != nil {
		return err
	}
	log.Debug("", "stdout", outbuf.String(), "stderr", errbuf.String())
	return os.Chdir(wd)
}
