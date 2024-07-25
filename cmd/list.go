package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
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
	entries, err := os.ReadDir(dotstashPath)
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		return errors.New("couldn't find any config files!")
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
		l.Item(newListItem(entryPath))
	}
	fmt.Println()
	fmt.Println(l)
	return nil
}

type listItem struct {
	title       string
	description string
	modules     []string
}

func (l listItem) String() string {
	titleStyle := lipgloss.NewStyle().Bold(true)
	descStyle := lipgloss.NewStyle().Italic(true)
	return lipgloss.JoinVertical(0,
		titleStyle.Render(l.title),
		descStyle.Render(l.description),
		renderBoxRow(l.modules, lipgloss.NormalBorder()),
	)
}

func renderBoxRow(items []string, border lipgloss.Border) string {
	var b strings.Builder

	b.WriteString(border.TopLeft)
	for i, item := range items {
		if i != 0 {
			b.WriteString(border.MiddleTop)
		}
		b.WriteString(strings.Repeat(border.Top, len(item)))
	}
	b.WriteString(border.TopRight + "\n")
	for _, item := range items {
		b.WriteString(border.Left)
		b.WriteString(item)
	}
	b.WriteString(border.Right + "\n")
	b.WriteString(border.BottomLeft)
	for i, item := range items {
		if i != 0 {
			b.WriteString(border.MiddleBottom)
		}
		b.WriteString(strings.Repeat(border.Top, len(item)))
	}
	b.WriteString(border.BottomRight)
	return b.String()
}

func newListItem(path string) string

func init() {
	rootCmd.AddCommand(listCmd)
	var err error
	dotstashPath, err = files.GetDotstashPath()
	if err != nil {
		panic("could not get dotstash path")
	}
}
