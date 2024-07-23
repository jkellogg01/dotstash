package cmd

import (
	"path"

	"github.com/charmbracelet/log"
	"github.com/jkellogg01/figure/files"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:  "remove name",
	Args: cobra.ExactArgs(1),
	RunE: removeFn,
}

func removeFn(cmd *cobra.Command, args []string) error {
	figRoot, err := files.GetFigurePath()
	if err != nil {
		return err
	}
	targetPath := path.Join(figRoot, args[0])
	log.Debug(targetPath)
	return nil
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
