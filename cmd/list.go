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
	t := table.New().Border(lipgloss.NormalBorder()).Headers("Primary", "Name", "Author", "Modules")
	for _, e := range entries {
		p := filepath.Join(dotstashPath, e.Name())
		meta, err := manifest.ReadManifest(p)
		if err != nil {
			log.Error("failed to get metadata, skipping...", "path", p, "error", err)
			continue
		}
		t.Row(fmt.Sprint(e.Name() == primary), e.Name(), meta.Author, targetsToNameList(meta.ExpandTargets()))
	}
	fmt.Println(t)
	log.Debug("", "primary", primary)
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

func init() {
	rootCmd.AddCommand(listCmd)
}
