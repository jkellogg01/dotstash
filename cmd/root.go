package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/jkellogg01/figure/files"
	"github.com/jkellogg01/figure/git"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "figure",
	Short: "An easy way to manage your configuration files",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	version, err := git.CheckGitInstalled()
	if err != nil {
		log.Fatal("error finding git installation!\n\tmake sure you have git installed; figure will not work without it.")
	} else if version == "" {
		log.Fatal("couldn't find a git installation!\n\tmake sure you have git installed; figure will not work without it.")
	}
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
	path, err := files.GetFigurePath()
	cobra.CheckErr(err)

	viper.AddConfigPath(path)
	viper.SetConfigType("toml")
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
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
