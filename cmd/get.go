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

var getCmd = &cobra.Command{
	Use:   "get [--branch=<branch name>] url",
	Short: "clones the repository at the specified url for use as a dotstash garden",
	RunE:  getFn,
	Args:  cobra.ExactArgs(1),
}

func getFn(cmd *cobra.Command, args []string) error {
	branch, err := cmd.Flags().GetString("branch")
	if err != nil {
		return err
	}
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
	alias, err := cmd.Flags().GetString("alias")
	if err != nil {
		return err
	}
	if alias != "" {
		name = alias
	}
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
	c := huh.NewConfirm().Title("Would you like to set this garden as your primary source for configuration files?")
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
	log.Infof("Successfully set %s as your primary configuration!", target)
	return nil
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().StringP("branch", "b", "", "specify a branch to download from")
	getCmd.Flags().StringP("alias", "a", "", "specify an alias for the created garden")
}
