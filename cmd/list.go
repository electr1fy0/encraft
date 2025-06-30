package cmd

import "github.com/spf13/cobra"


var listCmd = &cobra.Command{
	Use: "list",
	Short: "Show all entries",
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, args []string) {

	}
}
