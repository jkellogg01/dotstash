package cmd

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var (
	interactive bool
	atUrl       string
)

var getCmd = &cobra.Command{
	Use:   "get [url]",
	Short: "fetch a remote git repository to apply its configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Debug("get called", "args", args)
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
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "use an interactive prompt for downloading a config")
}
