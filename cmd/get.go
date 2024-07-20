package cmd

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var (
	interactive bool
	atUrl       string
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "A brief description of your command",
	RunE:  run,
	Args:  cobra.NoArgs,
}

func run(cmd *cobra.Command, args []string) error {
	log.Debug("get called", "args", args)
	return nil
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Use an interactive prompt for downloading a config")
	getCmd.Flags().StringVar(&atUrl, "url", "", "The url pointing to where the git repository is hosted")
	getCmd.Flags().StringP("protocol", "p", "ssh", "The protocol for git to use when cloning the repository")
	getCmd.Flags().StringP("hostname", "g", "github.com", "The hostname for the git host hosting the repository")
	getCmd.Flags().StringP("owner", "u", "", "The organization or user who owns the git repository")
	getCmd.Flags().StringP("repository", "r", "", "The name of the git repository")
	getCmd.MarkFlagsRequiredTogether("owner", "repository")
	getCmd.MarkFlagsMutuallyExclusive("interactive", "protocol", "url")
	getCmd.MarkFlagsMutuallyExclusive("interactive", "owner", "url")
	getCmd.MarkFlagsMutuallyExclusive("interactive", "owner", "url")
	getCmd.MarkFlagsOneRequired("interactive", "owner", "url")
}
