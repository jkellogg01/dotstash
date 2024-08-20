package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"strings"

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
	Short: "initialize a configuration garden and specify flowers to add",
	RunE:  makeFn,
}

func makeFn(cmd *cobra.Command, args []string) error {
	root, err := createConfigDir(dirName)
	if err != nil {
		return err
	}
	metadata := manifest.ConfigMetadata{
		Author: strings.TrimSpace(author),
	}
	log.Info("initializing a git repository in the new garden...")
	err = git.InitRepo(root)
	if err != nil {
		log.Error("failed to initialize git repository", "error", err)
	}
	if len(args) == 0 {
		metadata.EmitManifest(root)
		err = gardenInitialCommit(root)
		if err != nil {
			log.Warn(err)
		}
		return nil
	}
	wd, err := os.Getwd()
	if err != nil {
		log.Error("failed to get current workdir", "error", err)
	}
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
	log.Debug(metadata.ExpandTargets())
	err = metadata.EmitManifest(root)
	if err != nil {
		log.Error("failed to emit manifest")
		return err
	}
	log.Info("Garden successfully created!")
	// NOTE: at this point the creation of the garden is successful, so any errors should be 'quiet'
	err = gardenInitialCommit(root)
	if err != nil {
		log.Warn(err)
	}
	return nil
}

func gardenInitialCommit(dir string) error {
	err := os.Chdir(dir)
	if err != nil {
		log.Warn("failed to cd into garden directory", "error", err)
		return nil
	}
	log.Info("Creating initial commit...", "location", dir)
	err = git.Exec([]string{"add", "."})
	if err != nil {
		log.Warn("failed to execute 'git add .'", "error", err)
		return nil
	}
	err = git.Exec([]string{"commit", "--message=initial commit\r\n\r\nwith love from Dotstash"})
	if err != nil {
		log.Warn("failed to execute 'git commit --message=\"initial commit\r\n\r\nwith love from Dotstash\"", "error", err)
		return nil
	}
	return nil
}

// createConfigDir will append a path separator to the end of the path to the new directory.
func createConfigDir(name string) (string, error) {
	figRoot, err := files.GetDotstashPath()
	if err != nil {
		return "", err
	}
	newCfgPath := filepath.Join(figRoot, name)
	err = os.Mkdir(newCfgPath, 0o700)
	if errors.Is(err, fs.ErrExist) {
		log.Infof("directory '%s' already exists.", newCfgPath)
		i := 0
		for errors.Is(err, fs.ErrExist) {
			newCfgPath = filepath.Join(figRoot,
				fmt.Sprintf("%s_%03d", name, i))
			err = os.Mkdir(newCfgPath, 0o700)
		}
		if err != nil {
			log.Error("failed to create new garden under an alternate name", "error", err)
			return "", err
		}
	} else if err != nil {
		log.Error("failed to create new garden", "error", err)
		return "", err
	}
	log.Infof("successfully created a new garden at %s", newCfgPath)
	return newCfgPath, nil
}

func init() {
	rootCmd.AddCommand(makeCmd)
	makeCmd.Flags().StringVarP(&dirName, "name", "n", "dotstash", "the name of the garden to create")
	var defaultAuthorName string
	user, err := user.Current()
	if err == nil {
		defaultAuthorName = user.Username
	}
	makeCmd.Flags().StringVarP(&author, "author", "a", defaultAuthorName, "author name for the garden. defaults to blank if no username can be found")
	// TODO: add a flag for an interactive mode when there is an interactive mode to opt into
}
