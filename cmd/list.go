package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/charmbracelet/x/term"
	"github.com/electr1fy0/encraft/storage"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "Show all entries",
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, args []string) {
		exists, err := storage.VaultExists()

		if err != nil {
			fmt.Println("Error checking vault's presence: ", err)
			os.Exit(1)
		}

		if !exists {
			fmt.Println("No vault found. Initialize first")
			os.Exit(1)

		}
		fmt.Print("Enter master password: ")
		passwordBytes, err := term.ReadPassword(0)
		if err != nil {
			fmt.Println("Error reading master password")
			os.Exit(1)

		}
		password := string(passwordBytes)

		vault, err := storage.LoadVault(password)
		names := vault.ListEntries()

		// for i, entry := range vault.Entries {
		// 	fmt.Println(i+1, " ", entry.Name)
		// }
		sort.Strings(names)

		for i, name := range names {
			fmt.Println(i+1, ". ", name)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
