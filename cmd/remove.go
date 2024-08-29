package cmd

import (
	"os"
	"path"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"github.com/jkellogg01/dotstash/files"
	"github.com/jkellogg01/dotstash/manifest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var removeCmd = &cobra.Command{
	Use:   "remove name",
	Short: "remove the specified garden from the stored set of gardens, restoring the flowers therein to their proper locations if the selected garden was the primary garden",
	Args:  cobra.ExactArgs(1),
	RunE:  removeFn,
}

func removeFn(cmd *cobra.Command, args []string) error {
	figRoot, err := files.GetDotstashPath()
	if err != nil {
		return err
	}
	targetPath := path.Join(figRoot, args[0])
	norestore, err := cmd.Flags().GetBool("no-restore")
	if err != nil {
		return err
	}
	primary := viper.GetString("primary_config")
	if norestore || primary != args[0] {
		return os.RemoveAll(targetPath)
	}
	viper.Set("primary_config", "")
	err = viper.WriteConfig()
	if err != nil {
		log.Warn("failed to rewrite config", "error", err)
	}
	meta, err := manifest.ReadManifest(targetPath)
	if err != nil {
		return err
	}
	targets := meta.ExpandTargets()
	if os.Getenv("MODE") == "dev" {
		log.Debug("got yer targets here", "targets", targets)
		var didConfirm bool
		confirm := huh.NewConfirm().Value(&didConfirm).Title("dev mode: confirm delete")
		err = confirm.Run()
		if err != nil {
			return err
		}
		if !didConfirm {
			return nil
		}
		log.Debug("confirmed, starting deletion")
	}
	for _, t := range targets {
		log.Infof("substituting symlink at %s for %s", t.Dst, t.Src)
		err := files.Substitute(t.Src, t.Dst)
		if err != nil {
			return err
		}
	}
	log.Info("everything in its proper place. deleting directory...")
	return os.RemoveAll(targetPath)
}

func init() {
	rootCmd.AddCommand(removeCmd)

	removeCmd.Flags().Bool("no-restore", false, "do not restore configurations to their original locations, even if no other configuration is targeting that location")
}
