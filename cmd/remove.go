package cmd

import (
	"path"

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
	log.Debugf("%+v", meta)
	return nil
}

func init() {
	rootCmd.AddCommand(removeCmd)
	// TODO: this will probably also have an interactive mode at some point
}
