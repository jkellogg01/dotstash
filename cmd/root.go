package cmd

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/jkellogg01/figure/git"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "figure",
	Short: "An easy way to manage your configuration files",
}

func Execute() {
	rootCmd.SetErrPrefix(log.DefaultStyles().Levels[log.ErrorLevel].Render("ERRO"))
	version, err := git.CheckGitInstalled()
	if err != nil {
		log.Fatal("error finding git installation!\n\tmake sure you have git installed; figure will not work without it.")
	} else if version == "" {
		log.Fatal("couldn't find a git installation!\n\tmake sure you have git installed; figure will not work without it.")
	}
	log.Debug("checked and found git installation", "version", version)
	err = rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	path, err := os.UserConfigDir()
	cobra.CheckErr(err)

	viper.AddConfigPath(path)
	viper.SetConfigType("json")
	viper.SetConfigName("figure")

	viper.AutomaticEnv() // read in environment variables that match
	if os.Getenv("MODE") == "dev" {
		all := viper.AllSettings()
		log.Debugf("found %d config settings", len(all))
		for k, v := range all {
			log.Debug("%s = %v", k, v)
		}
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Infof("Using config file: %v", viper.ConfigFileUsed())
	}
}
