package main

import (
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/jkellogg01/figure/files"
	"github.com/jkellogg01/figure/git"
	"github.com/jkellogg01/figure/prompt"
)

func main() {
	log.SetLevel(log.DebugLevel)
	version, err := git.CheckGitInstalled()
	if err != nil {
		log.Fatalf("Failed to check git installation: %s", err)
	} else if version == "" {
		log.Fatal("Could not find git installation on this machine")
	}
	err = files.GotoFigureRoot()
	if err != nil {
		log.Fatal(err)
	}
}

func setupConfig() {
	// do nothing yet
}

func fetchConfig() {
	p := prompt.NewDefaultPrompter()
	host := p.TextPrompt("Enter your target git host", "github.com")
	owner := p.TextPrompt("Enter the owner of your target repo", "")
	repo := p.TextPrompt("Enter the name of your target repo", "dotfiles")
	ssh := p.BoolPrompt("Do you want to use ssh to clone this repo?", false)
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
