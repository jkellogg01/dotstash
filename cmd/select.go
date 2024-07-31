package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/jkellogg01/dotstash/manifest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var clobber, unlink bool

var selectCmd = &cobra.Command{
	Use:   "select <repo>",
	Short: "A brief description of your command",
	RunE:  selectFn,
	Args:  cobra.MinimumNArgs(1),
}

func selectFn(cmd *cobra.Command, args []string) error {
	primary := viper.GetString("primary_config")
	if primary == args[0] {
		log.Infof("%s is already your primary configuration!", primary)
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
		log.Error("failed to unlink config", "path", path, "error", err)
	}
	err = meta.Unlink()
	if err != nil {
		log.Error("failed to unlink config", "path", path, "error", err)
	}
}

func init() {
	rootCmd.AddCommand(selectCmd)

	selectCmd.Flags().BoolVarP(&clobber, "clobber", "c", false, "delete potentially non-symlink files when replacing them with configuration data from this repository")
	selectCmd.Flags().BoolVarP(&unlink, "unlink", "u", true, "unlink configuration files provided by the old repository")
}
