package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"path"

	"github.com/charmbracelet/log"
	"github.com/jkellogg01/figure/files"
	"github.com/jkellogg01/figure/manifest"
	"github.com/spf13/cobra"
)

var (
	dirName string
)

var makeCmd = &cobra.Command{
	Use:   "make [-n name] [file]...",
	Short: "set up a git repository and add config files to it",
	RunE:  makeFn,
}

func makeFn(cmd *cobra.Command, args []string) error {
	root, err := createConfigDir(dirName)
	if err != nil {
		return err
	}
	var metadata manifest.ConfigMetadata
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
		if err := linkSubstitute(oldPath, newPath); err != nil {
			return err
		}
	}
	metadata.EmitManifest(root)
	return nil
}

func linkSubstitute(oldPath, newPath string) error {
	backup, err := files.MakeTempFallback(oldPath)
	if err != nil {
		return err
	}
	defer backup.Close()
	backupInfo, err := backup.Stat()
	if err != nil {
		return err
	}
	backupName := backupInfo.Name()
	err = os.Rename(oldPath, newPath)
	if err != nil {
		log.Errorf("failed to move %s to %s. deleting backup and moving on...", oldPath, newPath)
		cleanupErr := os.RemoveAll(backupName)
		if cleanupErr != nil {
			log.Errorf("failed to clean up backup: %s", cleanupErr)
		}
		return err
	}
	err = errors.Join(
		os.Symlink(newPath, oldPath),
		os.Chmod(oldPath, 0o640),
	)
	if err == nil {
		err = os.RemoveAll(backupName)
		if err != nil {
			log.Errorf("backup not cleaned up: %s", err)
		}
		return nil
	}
	restoreBackupError := os.Rename(backupName, oldPath)
	err = errors.Join(err, restoreBackupError)
	if restoreBackupError != nil {
		log.Errorf("failed to restore %s from backup. backup is located at: %s", oldPath, backup)
	}
	return err
}

// createConfigDir will append a path separator to the end of the path to the new directory.
func createConfigDir(name string) (string, error) {
	figRoot, err := files.GetFigurePath()
	if err != nil {
		return "", err
	}
	newCfgPath := path.Join(figRoot, name)
	err = os.Mkdir(newCfgPath, 0o640)
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
		err = os.Mkdir(newCfgPath, 0o640)
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
	var defaultDirName string
	user, err := user.Current()
	if err != nil {
		defaultDirName = "dotfiles"
	} else {
		defaultDirName = user.Username
	}
	makeCmd.Flags().StringVarP(&dirName, "name", "n", defaultDirName, "the name of the config directory to create. Defaults to the username for the current user, or 'dotfiles' if no username is available")
	// TODO: add a flag for an interactive mode when there is an interactive mode to opt into
}
