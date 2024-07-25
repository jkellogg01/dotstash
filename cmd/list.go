package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/charmbracelet/log"
	"github.com/jkellogg01/dotstash/files"
	"github.com/jkellogg01/dotstash/manifest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	dotstashPath string
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "lists all stored configuration repositories.",
	RunE:  listFn,
	Args:  cobra.NoArgs,
}

func listFn(cmd *cobra.Command, args []string) error {
	// TODO: determine the primary set of configuration files
	entries, err := os.ReadDir(dotstashPath)
	if err != nil {
		return err
	}
	primary := viper.GetString("primary_config")
	if primary == "" {
		primary = entries[0].Name()
	}
	log.Debug("", "primary", primary)

	l := list.New().ItemStyleFunc(func(items list.Items, index int) lipgloss.Style {
		def := lipgloss.NewStyle().
			Padding(0, 1).MarginBottom(1).Border(lipgloss.NormalBorder(), false, false, false, true)
		if strings.Contains(items.At(index).Value(), primary) {
			highlight := lipgloss.Color("#F780E2")
			return def.Foreground(highlight).BorderForeground(highlight)
		}
		return def
	}).Enumerator(func(items list.Items, index int) string { return "" })
	for _, e := range entries {
		entryPath := filepath.Join(dotstashPath, e.Name())
		meta, err := manifest.ReadManifest(entryPath)
		if err != nil {
			continue
		}
		l.Item(listItem(e.Name(), "by "+meta.Author))
	}
	fmt.Println()
	fmt.Println(l)
	return nil
}

func listItem(title, description string) string {
	titleStyle := lipgloss.NewStyle().Bold(true)
	descStyle := lipgloss.NewStyle().Italic(true)
	return lipgloss.JoinVertical(0, titleStyle.Render(title), descStyle.Render(description))
}

func init() {
	rootCmd.AddCommand(listCmd)
	var err error
	dotstashPath, err = files.GetDotstashPath()
	if err != nil {
		panic("could not get dotstash path")
	}
}
