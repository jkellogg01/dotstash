package main

import (
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/jkellogg01/figure/prompt"
)

func main() {
	log.SetLevel(log.DebugLevel)
	haveGit, err := checkGitInstalled()
	if err != nil {
		log.Fatal("Failed to check git installation", "error", err)
	} else if !haveGit {
		log.Fatal("Could not find git installation on this machine")
	}
	p := prompt.NewDefaultPrompter()
	host := p.TextPrompt("Enter your target git host", "github.com")
	owner := p.TextPrompt("Enter the owner of your target repo", "")
	repo := p.TextPrompt("Enter the name of your target repo", "dotfiles")
	ssh := p.BoolPrompt("Do you want to us ssh to clone this repo?", false)
	var proto string
	if ssh {
		proto = "git@"
	} else {
		proto = "https://"
	}
	path := fmt.Sprintf("%s%s/%s/%s.git", proto, host, owner, repo)
	confirm := p.BoolPrompt(fmt.Sprintf("Does %s look right?", path), true)
	log.Info("", "path", path, "confirm", confirm)
}

func checkGitInstalled() (bool, error) {
	cmd := exec.Command("git", "--version")
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		return false, err
	}
	if err := cmd.Start(); err != nil {
		return false, err
	}
	data, err := io.ReadAll(pipe)
	if err != nil {
		return false, err
	}
	if err := cmd.Wait(); err != nil {
		return false, err
	}
	_, found := strings.CutPrefix(string(data), "git version ")
	return found, nil
}
