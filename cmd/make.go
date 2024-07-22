package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/user"

	"github.com/charmbracelet/log"
	"github.com/jkellogg01/figure/files"
	"github.com/spf13/cobra"
)

var (
	dirName string
)

var makeCmd = &cobra.Command{
	Use:   "make [-n name] [file]...",
	Short: "set up a git repository and add config files to it",
	Long: `make starts an interactive command prompt to set up a new git repository for your configuration files.
	you can optionally specify a starting dir; the prompt will start in $XDG_CONFIG_HOME by default.`,
	RunE: makeFn,
}

func makeFn(cmd *cobra.Command, args []string) error {
	root, err := createConfigDir(dirName)
	if err != nil {
		return err
	}
	if len(args) == 0 {
		return nil
	}
	wd, err := os.Getwd()
	if err != nil {
		log.Error("failed to get current workdir", "error", err)
	}
	for _, name := range args {
		log.Infof("adding %s...", name)
		info, err := os.Stat(name)
		if err != nil {
			return err
		}
		path := wd + string(os.PathSeparator) + info.Name()
		log.Debugf("%s is at path %s", name, path)
		err = os.Rename(path, root+name)
		if err != nil {
			return err
		}

	}
	return nil
}

// createConfigDir will append a path separator to the end of the path to the new directory.
func createConfigDir(name string) (string, error) {
	figRoot, err := files.GetFigurePath()
	if err != nil {
		return "", err
	}
	newCfgPath := files.PathAppend(figRoot, name)
	err = os.Mkdir(newCfgPath, os.ModeDir|777)
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
		err = os.Mkdir(newCfgPath, os.ModeDir|777)
		if err != nil {
			log.Error("failed to create new config dir", "error", err)
			return "", err
		}
	} else if err != nil {
		log.Error("failed to create new config dir", "error", err)
		return "", err
	}
	log.Infof("successfully created a new config dir at %s", newCfgPath)
	return newCfgPath + string(os.PathSeparator), nil
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
	// TODO add a flag for an interactive mode when there is an interactive mode to opt into
}
