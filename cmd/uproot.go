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

var uprootCmd = &cobra.Command{
	Use:     "uproot",
	Aliases: []string{"deplant"},
	Short:   "removes the specified flower from the primary garden, or a specified garden",
	RunE:    uprootFn,
	Args:    cobra.MinimumNArgs(1),
}

func uprootFn(cmd *cobra.Command, args []string) error {
	garden, err := cmd.Flags().GetString("garden")
	if err != nil {
		return err
	}
	repos, err := os.ReadDir(dotstashPath)
	if err != nil {
		return err
	}
	if garden == "" {
		garden = viper.GetString("primary_config")
		if garden == "" {
			return errors.New("no garden specified, and no primary garden")
		}
	}
	var target os.DirEntry
	for _, e := range repos {
		if e.Name() == garden {
			target = e
			break
		}
	}
	if target == nil {
		return fmt.Errorf("%s is not in your current list of gardens!", garden)
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
	rootCmd.AddCommand(uprootCmd)

	uprootCmd.Flags().StringP("garden", "g", "", "the garden to remove the flower(s) from")
}
