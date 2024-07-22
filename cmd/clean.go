package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"github.com/jkellogg01/figure/files"
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "[WARNING] DEVELOPMENT ONLY! clean will remove ALL of your stored figure directories.",
	Long: `[WARNING] THIS IS A DEVELOPER TOOL! IT IS NOT RECOMMENDED THAT YOU _EVER_ USE THIS COMMAND IN PRODUCTION [WARNING]
	clean will delete figure's root directory and everything inside it. This is intended to make repeated teasting of the program more convenient and is not recommended for figure users under any circumstances.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := files.GetFigurePath()
		if err != nil {
			return err
		}
		files, err := os.ReadDir(root)
		if err != nil {
			return err
		}
		fileCount := len(files)
		var confirm bool
		err = huh.NewConfirm().
			Title("WARNING: you are using a DESTRUCTIVE developer tool.").
			Description(fmt.Sprintf("This is not recommended for users under any circumstances.\nAre you sure you want to delete your figure directory and all %d entries it contains?", fileCount)).
			Affirmative("Delete my data").
			Negative("That seems bad").
			Value(&confirm).
			Run()
		if err != nil {
			return err
		}
		if !confirm {
			log.Info("clean cancelled.")
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
	rootCmd.AddCommand(cleanCmd)
}
