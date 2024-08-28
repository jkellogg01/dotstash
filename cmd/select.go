package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/jkellogg01/dotstash/manifest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var selectCmd = &cobra.Command{
	Use:   "select <garden>",
	Short: "select a garden from which to source your configuration files",
	RunE:  selectFn,
	Args:  cobra.MinimumNArgs(1),
}

func selectFn(cmd *cobra.Command, args []string) error {
	clobber, err := cmd.Flags().GetBool("clobber")
	if err != nil {
		return err
	}
	unlink, err := cmd.Flags().GetBool("unlink")
	if err != nil {
		return err
	}
	primary := viper.GetString("primary_config")
	if primary == args[0] {
		log.Infof("%s is already your primary garden!", primary)
		return nil
	}
	if unlink {
		unlinkRepo(filepath.Join(dotstashPath, primary))
	}
	path := filepath.Join(dotstashPath, args[0])
	meta, err := manifest.ReadManifest(path)
	if err != nil {
		return fmt.Errorf("error reading manifest at %s: %w", path, err)
	}
	meta.Link(clobber)
	viper.Set("primary_config", args[0])
	err = viper.WriteConfig()
	if err != nil {
		log.Errorf("Failed to write config: %v", err)
	}
	return nil
}

func unlinkRepo(path string) {
	meta, err := manifest.ReadManifest(path)
	if err != nil {
		log.Error("failed to unlink flower", "path", path, "error", err)
	}
	err = meta.Unlink()
	if err != nil {
		log.Error("failed to unlink flower", "path", path, "error", err)
	}
}

func init() {
	rootCmd.AddCommand(selectCmd)

	selectCmd.Flags().BoolP("clobber", "c", false, "delete potentially non-flower files when switching to the new garden")
	selectCmd.Flags().BoolP("unlink", "u", true, "unlink flowers provided by the old garden")
}
