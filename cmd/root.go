package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "encraft",
	Short: "manage your passwords from the command line",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}

}
