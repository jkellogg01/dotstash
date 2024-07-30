package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/jkellogg01/dotstash/manifest"
	"github.com/spf13/cobra"
)

var selectCmd = &cobra.Command{
	Use:   "select <repo>",
	Short: "A brief description of your command",
	RunE:  selectFn,
	Args:  cobra.MinimumNArgs(1),
}

func selectFn(cmd *cobra.Command, args []string) error {
	path := filepath.Join(dotstashPath, args[0])
	_, err := manifest.ReadManifest(path)
	if err != nil {
		return fmt.Errorf("error reading manifest at %s: %w", path, err)
	}
	// TODO: link from metadata targets
	return nil
}

func init() {
	rootCmd.AddCommand(selectCmd)
}
