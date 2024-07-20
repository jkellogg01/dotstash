package cmd

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var makeCmd = &cobra.Command{
	Use:   "make",
	Short: "set up a git repository and add config files to it",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Debug("make called", "args", args)
		return nil
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if interactive {
			return nil
		}
		return cobra.ExactArgs(1)(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(makeCmd)
}
