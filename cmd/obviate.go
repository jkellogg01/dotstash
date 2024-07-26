package cmd

import (
	"github.com/spf13/cobra"
)

var obviateCmd = &cobra.Command{
	Use:     "obviate",
	Aliases: []string{"rm-config"},
	Run:     obviateFn,
}

func obviateFn(cmd *cobra.Command, args []string) {

}

func init() {
	rootCmd.AddCommand(obviateCmd)
}
