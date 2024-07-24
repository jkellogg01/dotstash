package cmd

import (
	"os"
	"path"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"github.com/jkellogg01/dotstash/files"
	"github.com/jkellogg01/dotstash/manifest"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:  "remove name",
	Args: cobra.ExactArgs(1),
	RunE: removeFn,
}

func removeFn(cmd *cobra.Command, args []string) error {
	figRoot, err := files.GetDotstashPath()
	if err != nil {
		return err
	}
	targetPath := path.Join(figRoot, args[0])
	log.Debug(targetPath)
	// TODO: when storing and swapping between multiple configs is supported,
	// this should check which config is set as 'primary' and/or if it is being
	// referenced for any config targets. for now we will just assume the named
	// directory is primary
	meta, err := manifest.ReadManifest(targetPath)
	if err != nil {
		return err
	}
	targets, err := meta.Targets()
	if err != nil {
		return err
	}
	if os.Getenv("MODE") == "dev" {
		for _, t := range targets {
			log.Debugf("%s => %s", t.Src, t.Dst)
		}
		var didConfirm bool
		confirm := huh.NewConfirm().Value(&didConfirm).Title("dev mode: confirm delete")
		err = confirm.Run()
		if err != nil {
			return err
		}
		if !didConfirm {
			return nil
		}
	}
	log.Debug("confirmed, starting deletion")
	return nil
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
