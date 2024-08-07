package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/charmbracelet/log"
	"github.com/jkellogg01/dotstash/manifest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "lists all stored configuration gardens.",
	RunE:  listFn,
	Args:  cobra.NoArgs,
}

func listFn(cmd *cobra.Command, args []string) error {
	var (
		indigo  = lipgloss.AdaptiveColor{Light: "#5A56E0", Dark: "#7571F9"}
		fuchsia = lipgloss.Color("#F780E2")
	)
	entries, err := os.ReadDir(dotstashPath)
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		log.Error("couldn't find any config files!")
		return nil
	}
	primary := viper.GetString("primary_config")
	if primary == "" {
		primary = entries[0].Name()
		viper.Set("primary_config", primary)
		// HACK: we shouldn't ever actually need to do this, it's mostly here for testing
		err := viper.WriteConfig()
		if err != nil {
			log.Error("failed to write to config", "error", err)
		}
	}
	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(indigo)).
		BorderColumn(false).
		// BorderHeader(false).
		Headers("Primary", "Name", "Author", "Modules")
	var primaryRow int
	currentRow := 1
	for _, e := range entries {
		p := filepath.Join(dotstashPath, e.Name())
		meta, err := manifest.ReadManifest(p)
		if err != nil {
			log.Error("failed to get metadata", "path", p, "error", err)
			t.Row("n/a", e.Name(), "n/a", "n/a")
			continue
		}
		isPrimary := e.Name() == primary
		if isPrimary {
			primaryRow = currentRow
		}
		t.Row(cSprint(isPrimary, "y", "n"), e.Name(), meta.Author, targetsToNameList(meta.ExpandTargets()))
		currentRow++
	}
	log.Debug("", "primary", primary)
	log.Debug("", "primary", primaryRow)
	t = t.StyleFunc(func(row, col int) lipgloss.Style {
		style := lipgloss.NewStyle().Padding(0, 1)
		if row == 0 {
			style = style.Bold(true).Foreground(indigo)
		}
		if primaryRow != 0 && row == primaryRow {
			style = style.Italic(true).Foreground(fuchsia)
		}
		return style
	})

	fmt.Println(t)
	return nil
}

func targetsToNameList(targets []manifest.ConfigTarget) string {
	var s strings.Builder
	for i, t := range targets {
		name := filepath.Base(t.Src)
		if i > 0 {
			s.WriteString(", ")
		}
		s.WriteString(name)
	}
	return s.String()
}

func cSprint(cond bool, ifTrue, ifFalse string) string {
	if cond {
		return ifTrue
	}
	return ifFalse
}

func init() {
	rootCmd.AddCommand(listCmd)
}
