package cmd

import (
	"errors"
	"net/url"
	"path"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/jkellogg01/dotstash/git"
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
	return nil
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().StringVarP(&branch, "branch", "b", "", "specify a branch to download from")
}
