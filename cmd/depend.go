package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/jkellogg01/dotstash/files"
	"github.com/jkellogg01/dotstash/manifest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var repoName string

var plantCmd = &cobra.Command{
	Use:  "plant [-r repository] path...",
	RunE: plantFunc,
	Args: cobra.MinimumNArgs(1),
}

func plantFunc(cmd *cobra.Command, args []string) error {
	repos, err := os.ReadDir(dotstashPath)
	if err != nil {
		return err
	}
	var target os.DirEntry
	if repoName == "" {
		repoName = viper.GetString("primary_config")
		if repoName == "" {
			return errors.New("no repository specified, and no primary repository")
		}
	}
	for _, e := range repos {
		if e.Name() == repoName {
			target = e
			break
		}
	}
	if target == nil {
		return fmt.Errorf("%s is not in your current list of repositories!", repoName)
	}
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	root := filepath.Join(dotstashPath, target.Name())
	metadata, err := manifest.ReadManifest(root)
	log.Debug(args)
	for _, name := range args {
		log.Infof("adding %s...", name)
		name := filepath.Clean(name)
		var oldPath, newPath string
		if filepath.IsAbs(name) {
			dir, name := filepath.Split(name)
			oldPath = filepath.Join(dir, name)
			newPath = filepath.Join(root, name)
		} else {
			oldPath = filepath.Join(wd, name)
			newPath = filepath.Join(root, name)
		}
		log.Debug("got the following paths", "old", oldPath, "new", newPath)
		metadata.AppendTarget(newPath, oldPath)
		if err := files.SubstituteForSymlink(oldPath, newPath); err != nil {
			return err
		}
	}
	metadata.EmitManifest(root)
	return nil
}

func init() {
	rootCmd.AddCommand(plantCmd)

	plantCmd.Flags().StringVarP(&repoName, "repository", "r", "", "the repository to add the flower to")
}
