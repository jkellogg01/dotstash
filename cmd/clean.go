package cmd

import (
	"os"

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
		var confirm bool
		err := huh.NewConfirm().
			Title("WARNING: you are using a DESTRUCTIVE developer tool.").
			Description("This is not recommended for users under any circumstances.\nAre you sure you want to delete your figure directory and ALL of its contents?").
			Affirmative("I want to delete all of my data").
			Negative("That seems bad").
			Value(&confirm).
			Run()
		if err != nil {
			return err
		}
		if !confirm {
			log.Info("clean cancelled.")
		}
		root, err := files.GetFigurePath()
		if err != nil {
			return err
		}
		return os.RemoveAll(root)
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}
