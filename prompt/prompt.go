package prompt

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/log"
)

type Prompter struct {
	scanner *bufio.Scanner
	writer  io.Writer
}

func NewPrompter(s *bufio.Scanner, w io.Writer) *Prompter {
	return &Prompter{s, w}
}

func NewDefaultPrompter() *Prompter {
	return NewPrompter(bufio.NewScanner(os.Stdin), os.Stdout)
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
	fmt.Fprint(p.writer, q)
	if d != "" {
		fmt.Fprintf(p.writer, " (%s)", d)
	}
	fmt.Fprint(p.writer, "\n> ")
	if !p.scanner.Scan() {
		if err := p.scanner.Err(); err != nil {
			log.Error(err)
		}
		return ""
	}
	return p.scanner.Text()
}
