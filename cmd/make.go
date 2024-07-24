package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"path"

	"github.com/charmbracelet/log"
	"github.com/jkellogg01/dotstash/files"
	"github.com/jkellogg01/dotstash/git"
	"github.com/jkellogg01/dotstash/manifest"
	"github.com/spf13/cobra"
)

var (
	dirName string
	author  string
)

var makeCmd = &cobra.Command{
	Use:   "make [flags] [file]...",
	Short: "set up a git repository and add config files to it",
	RunE:  makeFn,
}

func makeFn(cmd *cobra.Command, args []string) error {
	root, err := createConfigDir(dirName)
	if err != nil {
		return err
	}
	metadata := manifest.ConfigMetadata{
		Author: author,
	}
	log.Info("initializing a git repository in the new directory...")
	err = git.InitRepo(root)
	if err != nil {
		log.Error("failed to initialize git repository", "error", err)
	}
	if len(args) == 0 {
		metadata.EmitManifest(root)
		return nil
	}
	wd, err := os.Getwd()
	if err != nil {
		log.Error("failed to get current workdir", "error", err)
	}
	for _, name := range args {
		log.Infof("adding %s...", name)
		var oldPath, newPath string
		if path.IsAbs(name) {
			dir, name := path.Split(name)
			oldPath = path.Join(dir, name)
			newPath = path.Join(root, name)
		} else {
			oldPath = path.Join(wd, name)
			newPath = path.Join(root, name)
		}
		metadata.AppendTarget(newPath, oldPath)
		log.Debug("got the following paths", "old", oldPath, "new", newPath)
		if err := files.LinkSubstitute(oldPath, newPath); err != nil {
			return err
		}
	}
	metadata.EmitManifest(root)
	return nil
}

// createConfigDir will append a path separator to the end of the path to the new directory.
func createConfigDir(name string) (string, error) {
	figRoot, err := files.GetDotstashPath()
	if err != nil {
		return "", err
	}
	newCfgPath := path.Join(figRoot, name)
	err = os.Mkdir(newCfgPath, 0o700)
	if errors.Is(err, fs.ErrExist) {
		log.Infof("directory '%s' already exists. backing up and replacing...", newCfgPath)
		i := 0
		backupPath := fmt.Sprintf("%s_backup_%04d", newCfgPath, i)
		err = os.Rename(newCfgPath, backupPath)
		for errors.Is(err, fs.ErrExist) {
			i++
			backupPath = fmt.Sprintf("%s_backup_%04d", newCfgPath, i)
			err = os.Rename(newCfgPath, backupPath)
		}
		if err != nil {
			log.Error("failed to create backup", "error", err)
			return "", err
		}
		err = os.Mkdir(newCfgPath, 0o700)
		if err != nil {
			log.Error("failed to create new config dir", "error", err)
			return "", err
		}
	} else if err != nil {
		log.Error("failed to create new config dir", "error", err)
		return "", err
	}
	log.Infof("successfully created a new config dir at %s", newCfgPath)
	return newCfgPath, nil
}

func init() {
	rootCmd.AddCommand(makeCmd)
	makeCmd.Flags().StringVarP(&dirName, "name", "n", "dotstash", "the name of the config directory to create")
	var defaultAuthorName string
	user, err := user.Current()
	if err == nil {
		defaultAuthorName = user.Username
	}
	makeCmd.Flags().StringVarP(&author, "author", "a", defaultAuthorName, "author name for the repository. defaults to blank if no usernmae can be found")
	// TODO: add a flag for an interactive mode when there is an interactive mode to opt into
}
