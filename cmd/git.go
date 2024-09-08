package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// gitCmd represents the git command
var gitCmd = &cobra.Command{
	Use:                   "git [--garden=<garden name>] command...",
	Short:                 "a wrapper around git commands, which executes them in the primary or specified garden.",
	RunE:                  gitFn,
	Args:                  cobra.MinimumNArgs(1),
	DisableFlagParsing:    true,
	DisableFlagsInUseLine: true,
}

func gitFn(cmd *cobra.Command, args []string) error {
	// HACK: this is horrible; there must be a better way
	var targetGarden string
	switch {
	case args[0] == "--garden" || args[0] == "-g":
		targetGarden = args[1]
		args = args[2:]
		err := cmd.Args(cmd, args)
		if err != nil {
			return err
		}
	case strings.HasPrefix(args[0], "--garden="):
		t, _ := strings.CutPrefix(args[0], "--garden=")
		targetGarden = t
		args = args[1:]
		err := cmd.Args(cmd, args)
		if err != nil {
			return err
		}
	case args[0] == "--help" || args[0] == "-h":
		return cmd.Help()
	}
	log.Debug("remaining arguments", "args", args)
	repos, err := os.ReadDir(dotstashPath)
	if err != nil {
		return err
	}
	var target os.DirEntry
	if targetGarden == "" {
		n := viper.GetString("primary_config")
		log.Debugf("viper primary_config: %s", n)
		if n == "" {
			return errors.New("no garden specified, and no primary garden")
		}
		targetGarden = n
	}
	for _, e := range repos {
		if e.Name() == targetGarden {
			target = e
			break
		}
	}
	if target == nil {
		return fmt.Errorf("%s is not in your current list of gardens!", targetGarden)
	}
	targetPath := filepath.Join(dotstashPath, targetGarden)
	c := exec.Command(
		"git",
		args...,
	)
	c.Dir = targetPath
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	return c.Run()
}

func init() {
	rootCmd.AddCommand(gitCmd)

	gitCmd.Flags().StringP("garden", "g", "", "the garden in which to execute the specified git command")
}
