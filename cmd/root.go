package cmd

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/jkellogg01/dotstash/files"
	"github.com/jkellogg01/dotstash/git"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	dotstashPath string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dotstash",
	Short: "tending, harvesting, and sharing your configuration files has never been easier!",
	Long: `dotstash is a configuration management program which helps you manage your configuration files in a variety of ways.
- simple, readable metadata files to make setting up a new configuration, like, super easy!
- one-command process for downloading and trying out someone else's set of configuration files!
- easily add and remove configuration files for different applications!

For maximal clarity, dotstash uses garden terminology to specify the scope of file you're manipulating at any given time. A complete set of configuration files for a variety of applications is a _garden_. a _garden_ contains many _flowers_, each of which is the complete set of configuration files for a particular application.

This terminology may feel silly at first, but after working on this project for as long as I have, I can promise you that you will be so glad that you don't have to think about what exactly i mean by a 'repository' every time you RTFM.

Let's get more specific about terminology, as this will tell us a lot about the usage of dotstash:
- A **garden** maintains a _single branch_ of a git repository full of configuration files. A graden can be:
  - 'make'd to create a new garden, or 'sow'ed, if you please.
  - 'remove'd to, well, remove a garden, as well as restore the contained _flowers_ to their original homes if they are currently being used. This action can also be performed with the 'reap' command if you so desire.
- A **flower** is the specific set of configuration files for a particular application, such as the one you may find at '~/.config/sway'. A flower can be:
  - 'plant'ed in a _garden_ in order to add that set of configuration files to the garen
  - 'uproot'ed (or 'deplant'ed) from a _garden_ in order to remove that set of configuration files from the _garden's_ area of responsibility; restoring the files to their original locations if applicable.
- Finally, the set of commands that apply to all of your gardens in some way or another:
  - 'list' will present a list of your _gardens_ and the _flowers_ in their care.
  - 'select' will choose a primary _garden_ which will be the central source of configuration files. Future plans include selecting primary _flowers_ from non-primary _gardens_, but this functionality is still in the works.
  - 'torch' will **destructively remove all of your flowers without attempting to restore the configuration files therein to their original locations!** I cannot stress enough that this is a **destructive command** that i wrote for my own convenience as a developer and that you **should not use it, probably ever.** Cool? Cool.`,
}

func Execute() {
	rootCmd.SetErrPrefix(log.DefaultStyles().Levels[log.ErrorLevel].Render("ERRO"))
	version, err := git.CheckGitInstalled()
	if err != nil {
		log.Fatal("error finding git installation!\n\tmake sure you have git installed; dotstash will not work without it.")
	} else if version == "" {
		log.Fatal("couldn't find a git installation!\n\tmake sure you have git installed; dotstash will not work without it.")
	}
	log.Debug("checked and found git installation", "version", version)
	err = rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	var err error
	dotstashPath, err = files.GetDotstashPath()
	if err != nil {
		panic("could not get dotstash path")
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	path, err := os.UserConfigDir()
	cobra.CheckErr(err)

	viper.AddConfigPath(path)
	viper.SetConfigType("json")
	viper.SetConfigName("dotstash")

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Debugf("Using config file: %v", viper.ConfigFileUsed())
	} else {
		viper.SafeWriteConfig()
	}

	if os.Getenv("MODE") == "dev" {
		all := viper.AllSettings()
		log.Debugf("found %d config settings", len(all))
		for k, v := range all {
			log.Debugf("%s = %v", k, v)
		}
	}
}
