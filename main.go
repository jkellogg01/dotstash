package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
)

func main() {
	log.SetLevel(log.DebugLevel)
	var formData struct {
		Host       string
		Owner      string
		Repository string
	}
	f := huh.NewForm(huh.NewGroup(
		huh.NewInput().
			Title("Target git host:").
			Placeholder("github.com").
			Suggestions([]string{
				"github.com",
				"codeberg.org",
				"gitlab.com",
			}).
			Value(&formData.Host),
		huh.NewInput().
			Title("Code owner:").
			Placeholder("probably you?").
			Value(&formData.Owner),
		huh.NewInput().
			Title("Repository name:").
			Placeholder(".dotfiles").
			Value(&formData.Repository),
	))
	if err := f.Run(); err != nil {
		log.Fatal(err)
	}
	path := fmt.Sprintf(
		"https://%s.git",
		strings.Join([]string{
			formData.Host,
			formData.Owner,
			formData.Repository},
			"/"),
	)
	var looksGood bool
	huh.NewConfirm().
		Title("Your dotfiles should be here, right?").
		Description(path).
		Affirmative("That's right!").
		Negative("No...").
		Value(&looksGood).
		Run()
	log.Info(path, "looks good", looksGood)
}
