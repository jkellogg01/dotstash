package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/log"
)

type prompter struct {
	scanner *bufio.Scanner
}

func main() {
	log.SetLevel(log.DebugLevel)
	p := NewDefaultPrompter()
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

func (p *prompter) TextPrompt(q, d string) string {
	var defaultText string
	if d != "" {
		defaultText = fmt.Sprintf("default: %s", d)
	}
	resp := p.prompt(q, defaultText)
	if resp == "" {
		return d
	}
	return resp
}

func (p *prompter) BoolPrompt(q string, d bool) bool {
	var defaultText string
	if d {
		defaultText = "Y/n"
	} else {
		defaultText = "y/N"
	}
	resp := p.prompt(q, defaultText)
	if len(resp) == 0 {
		return d
	}
	switch strings.ToLower(resp)[0] {
	case 'y':
		return true
	case 'n':
		return false
	default:
		return d
	}
}

func (p *prompter) prompt(q, d string) string {
	fmt.Print(q)
	if d != "" {
		fmt.Printf(" (%s)", d)
	}
	fmt.Print("\n> ")
	if !p.scanner.Scan() {
		if err := p.scanner.Err(); err != nil {
			log.Error(err)
		}
		return ""
	}
	return p.scanner.Text()
}

func NewDefaultPrompter() *prompter {
	return &prompter{
		bufio.NewScanner(os.Stdin),
	}
}
