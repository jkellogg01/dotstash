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

var plantCmd = &cobra.Command{
	Use:   "plant path...",
	Short: "adds the specified flower to the primary garden, or a specified garden",
	RunE:  plantFunc,
	Args:  cobra.MinimumNArgs(1),
}

func plantFunc(cmd *cobra.Command, args []string) error {
	garden, err := cmd.Flags().GetString("garden")
	if err != nil {
		return err
	}
	repos, err := os.ReadDir(dotstashPath)
	if err != nil {
		return err
	}
	var target os.DirEntry
	if garden == "" {
		garden = viper.GetString("primary_config")
		if garden == "" {
			return errors.New("no garden specified, and no primary garden")
		}
	}
	for _, e := range repos {
		if e.Name() == garden {
			target = e
			break
		}
	}
	if target == nil {
		return fmt.Errorf("%s is not in your current list of gardens!", garden)
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

	plantCmd.Flags().StringP("garden", "g", "", "the garden to add the flower to")
}
