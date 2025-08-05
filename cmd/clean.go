package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/x/term"
	"github.com/electr1fy0/encraft/storage"
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "delete all items from the vault",
	Run: func(cmd *cobra.Command, args []string) {
		exists, err := storage.VaultExists()

		if err != nil {
			fmt.Println("Error checking vault's presence: ", err)
			os.Exit(1)
		}

		if !exists {
			fmt.Println("No vault inside home directory. Initialize it first")
			os.Exit(1)

		}

		fmt.Print("Enter master password: ")
		passwordBytes, err := term.ReadPassword(0)
		password := string(passwordBytes)

		vault, err := storage.LoadVault(password)
		if err != nil {
			fmt.Println("Error loading the vault")
			os.Exit(1)
		}

		for _, entry := range vault.Entries {
			vault.DeleteEntry(entry.Name)
		}

		err = storage.SaveVault(vault, password)

		fmt.Println("\nVault is now empty")
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}
