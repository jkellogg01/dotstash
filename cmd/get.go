package cmd

import (
	"errors"
	"io/fs"
	"net/url"
	"path"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"github.com/jkellogg01/dotstash/git"
	"github.com/jkellogg01/dotstash/manifest"
	"github.com/spf13/cobra"
)

var branch string

var getCmd = &cobra.Command{
	Use:   "get [--branch branchname] url",
	Short: "clones the repository at the specified url for use as a dotstash configuration",
	RunE:  getFn,
	Args:  cobra.ExactArgs(1),
}

func getFn(cmd *cobra.Command, args []string) error {
	if strings.HasPrefix(args[0], "git@") {
		// TODO: might have to make some changes here to support ssh cloning, so that people can still make commits to their own repositories.
		// the solution for now will just be to manually ssh clone dotfile repositories.
		return errors.New("please use the http url to the repository!")
	}
	t, _ := strings.CutSuffix(args[0], ".git")
	src, err := url.Parse(t)
	if err != nil {
		return err
	}
	name := path.Base(src.Path)
	// TODO: once there is an alias flag available, this should use that if it's non-blank
	target := filepath.Join(dotstashPath, name)
	err = git.Download(src.String(), branch, target)
	if err != nil {
		return err
	}
	log.Infof("successfully cloned %s into %s!", src.String(), target)
	meta, err := manifest.ReadManifest(target)
	if errors.Is(err, fs.ErrNotExist) {
		log.Warn("the downloaded repository does not contain a manifest.json; it will need one before it can be used with dotstash!")
		return nil
	}
	c := huh.NewConfirm().Title("Would you like to set this repo as your primary source for configuration files?")
	err = c.Run()
	if err != nil {
		log.Error("failed to run confirm prompt, exiting", "error", err)
		return nil
	}
	setPrimary, ok := c.GetValue().(bool)
	if !ok {
		panic("confirm field did not return a bool value")
	}
	if !setPrimary {
		log.Info("Job's done!")
		return nil
	}
	meta.Link(false)
	return nil
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().StringVarP(&branch, "branch", "b", "", "specify a branch to download from")
}
