package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"github.com/jkellogg01/dotstash/files"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var torchCmd = &cobra.Command{
	Use:   "torch",
	Short: "[WARNING] DEVELOPMENT ONLY! torch will remove ALL of your stored dotstash directories.",
	Long: `[WARNING] THIS IS A DEVELOPER TOOL! IT IS NOT RECOMMENDED THAT YOU _EVER_ USE THIS COMMAND IN PRODUCTION [WARNING]
	torch will delete dotstash's root directory and everything inside it. This is intended to make repeated teasting of the program more convenient and is not recommended for dotstash users under any circumstances.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := files.GetDotstashPath()
		if err != nil {
			return err
		}
		files, err := os.ReadDir(root)
		if err != nil {
			return err
		}
		fileCount := len(files)
		c := huh.NewConfirm().
			Title("WARNING: you are using a DESTRUCTIVE developer tool.").
			Description(fmt.Sprintf("This is not recommended for users under any circumstances.\nAre you sure you want to delete your dotstash directory and all %d entries it contains?", fileCount)).
			Affirmative("Delete my data").
			Negative("That seems bad").
			WithTheme(huh.ThemeBase())
		err = c.Run()
		if err != nil {
			return err
		}
		confirm, ok := c.GetValue().(bool)
		if !ok {
			panic("failed to get confirm value, get me out of here man")
		}
		if !confirm {
			log.Info("torch cancelled.")
			return nil
		}
		viper.Set("primary_config", "")
		err = viper.WriteConfig()
		if err != nil {
			log.Warn("failed to write config", "error", err)
		}
		for _, file := range files {
			name := path.Join(root, file.Name())
			log.Debugf("Removing dir %s", name)
			err := os.RemoveAll(name)
			if err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(torchCmd)
}
