package cmd

import (
	"github.com/spf13/cobra"
)

var (
	// Used for flags.
	cfgFile    string
	cfgVerbose bool

	rootCmd = &cobra.Command{
		Use:   "cleanABag",
		Short: "A tool to remove old articles from wallabag",
		Long:  `cleanABag is a CLI tool for removing all articles from wallabag.`,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "credentials", "c", "", "config file with credentials to connect to wallabag (default is $HOME/.config/cleanABag/credentials.json) - Full path is needed.")
	rootCmd.PersistentFlags().BoolVarP(&cfgVerbose, "verbose", "v", false, "Verbose mode")
}
