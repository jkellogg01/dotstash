package cmd

import (
	"context"
	"errors"
	"io/fs"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"github.com/jkellogg01/figure/files"
	"github.com/spf13/cobra"
)

var makeForm = huh.NewForm(huh.NewGroup())

var makeCmd = &cobra.Command{
	Use:   "make [dir]",
	Short: "set up a git repository and add config files to it",
	Long: `make starts an interactive command prompt to set up a new git repository for your configuration files.
	you can optionally specify a starting dir; the prompt will start in $XDG_CONFIG_HOME by default.`,
	RunE: makeFn,
	Args: cobra.MaximumNArgs(1),
}

func makeFn(cmd *cobra.Command, args []string) error {
	path, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	if len(args) == 1 {
		argPath := args[0]
		log.Debugf("attempting to start command prompt in %s", argPath)
		v, err := files.DirExists(argPath)
		if err != nil {
			return err
		}
		if v {
			path = argPath
		}
	}
	err = makeForm.RunWithContext(context.WithValue(
		context.Background(),
		"root",
		path,
	))
	if err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(makeCmd)
}
