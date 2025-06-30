package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/x/term"
	"github.com/electr1fy0/encraft/storage"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get [name]",
	Short: "Get your password for a name",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		entryName := args[0]

		exists, err := storage.VaultExists()
		if err != nil {
			fmt.Println("Error checking vault's presence: ", err)
			os.Exit(1)
		}

		if !exists {
			fmt.Println("Vault does not exist. Initialize first")
			os.Exit(1)
		}

		fmt.Print("Enter master pass: ")
		passwordBytes, err := term.ReadPassword(0)

		if err != nil {
			fmt.Println("Error reading master password: ", err)
			os.Exit(1)
		}
		password := string(passwordBytes)
		vault, err := storage.LoadVault(password)

		entry, exists := vault.GetEntry(entryName)
		if !exists {
			fmt.Println("Entry does not exist or we probably corrupted it. Too bad either way")
			os.Exit(1)
		}

		fmt.Printf("----%s----\n", entry.Name)

		fmt.Printf("----%s----\n", entry.Password)
		fmt.Println("Created at: ", entry.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Println("Updated at: ", entry.UpdatedAt.Format("2006-01-02 15:04:05"))
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
