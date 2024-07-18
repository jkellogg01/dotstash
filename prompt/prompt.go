package prompt

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/log"
)

type Prompter struct {
	scnr *bufio.Scanner
}

func NewPrompter(s *bufio.Scanner) *Prompter {
	return &Prompter{s}
}

func NewDefaultPrompter() *Prompter {
	return NewPrompter(bufio.NewScanner(os.Stdin))
}

func (p *Prompter) TextPrompt(q, d string) string {
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

func (p *Prompter) BoolPrompt(q string, d bool) bool {
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

func (p *Prompter) prompt(q, d string) string {
	fmt.Print(q)
	if d != "" {
		fmt.Printf(" (%s)", d)
	}
	fmt.Print("\n> ")
	if !p.scnr.Scan() {
		if err := p.scnr.Err(); err != nil {
			log.Error(err)
		}
		return ""
	}
	return p.scnr.Text()
}
