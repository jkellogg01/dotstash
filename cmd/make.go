package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/charmbracelet/log"
	"github.com/jkellogg01/figure/files"
	"github.com/spf13/cobra"
)

var makeCmd = &cobra.Command{
	Use:   "make name [file]...",
	Short: "set up a git repository and add config files to it",
	Long: `make starts an interactive command prompt to set up a new git repository for your configuration files.
	you can optionally specify a starting dir; the prompt will start in $XDG_CONFIG_HOME by default.`,
	RunE: makeFn,
	Args: cobra.MinimumNArgs(1),
}

func makeFn(cmd *cobra.Command, args []string) error {
	figRoot, err := files.GetFigurePath()
	if err != nil {
		return err
	}
	newCfgPath := figRoot + string(os.PathSeparator) + args[0]
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
			return err
		}
		err = os.Mkdir(newCfgPath, os.ModeDir|777)
		if err != nil {
			log.Error("failed to create new config dir", "error", err)
			return err
		}
	} else if err != nil {
		log.Error("failed to create new config dir", "error", err)
		return err
	}
	log.Debugf("successfully created a new config dir at %s", newCfgPath)
	return nil
}

func init() {
	rootCmd.AddCommand(makeCmd)
}
