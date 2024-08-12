package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var targetGarden string

// gitCmd represents the git command
var gitCmd = &cobra.Command{
	Use:   "git [--garden=<garden name>] command...",
	Short: "a wrapper around git commands, which executes them in the primary or specified garden.",
	RunE:  gitFn,
	Args:  cobra.MinimumNArgs(1),
}

func gitFn(cmd *cobra.Command, args []string) error {
	repos, err := os.ReadDir(dotstashPath)
	if err != nil {
		return err
	}
	var target os.DirEntry
	if repoName == "" {
		if n := viper.GetString("primary_config"); n != "" {
			repoName = n
		}
		return errors.New("no garden specified, and no primary garden")
	}
	for _, e := range repos {
		if e.Name() == repoName {
			target = e
			break
		}
	}
	if target == nil {
		return fmt.Errorf("%s is not in your current list of gardens!", repoName)
	}
	// TODO: cd into the specified garden and execute the provided git command
	return nil
}

func init() {
	rootCmd.AddCommand(gitCmd)

	gitCmd.Flags().StringVarP(&targetGarden, "garden", "g", "", "the garden in which to execute the specified git command")
}
