package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/charmbracelet/log"
	"github.com/jkellogg01/dotstash/files"
	"github.com/jkellogg01/dotstash/manifest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var obviateCmd = &cobra.Command{
	Use:     "obviate",
	Aliases: []string{"rm-config"},
	RunE:    obviateFn,
	Args:    cobra.MinimumNArgs(1),
}

func obviateFn(cmd *cobra.Command, args []string) error {
	repos, err := os.ReadDir(dotstashPath)
	if err != nil {
		return err
	}
	if repoName == "" {
		repoName = viper.GetString("primary_config")
		if repoName == "" {
			return errors.New("no repository specified, and no primary repository")
		}
	}
	var target os.DirEntry
	for _, e := range repos {
		if e.Name() == repoName {
			target = e
			break
		}
	}
	if target == nil {
		return fmt.Errorf("%s is not in your current list of repositories!", repoName)
	}
	targetPath := filepath.Join(dotstashPath, target.Name())
	metadata, err := manifest.ReadManifest(targetPath)
	if err != nil {
		return err
	}
	for _, t := range metadata.ExpandTargets() {
		basename := filepath.Base(t.Src)
		if !slices.Contains(args, basename) {
			continue
		}
		err := files.Substitute(t.Src, t.Dst)
		if err != nil {
			log.Error("substitution failed", "target", t, "error", err)
		}
		metadata.RemoveTarget(basename)
	}
	return metadata.EmitManifest(targetPath)
}

func init() {
	rootCmd.AddCommand(obviateCmd)

	obviateCmd.Flags().StringVarP(&repoName, "repository", "r", "", "the repository to remove the config from")
}
