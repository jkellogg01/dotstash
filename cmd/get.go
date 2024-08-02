package cmd

import (
	"strings"

	"github.com/spf13/cobra"
)

var branch string

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "A brief description of your command",
	RunE:  getFn,
	Args:  cobra.ExactArgs(1),
}

func getFn(cmd *cobra.Command, args []string) error {
	var info strings.Builder
	info.WriteString(args[0])
	if branch != "" {
		info.WriteString(" on branch " + branch)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().StringVarP(&branch, "branch", "b", "", "specify a branch to download from")
}
