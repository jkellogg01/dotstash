package main

import (
	"bufio"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/jkellogg01/figure/prompt"
)

type prompter struct {
	scanner *bufio.Scanner
}

func main() {
	log.SetLevel(log.DebugLevel)
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
