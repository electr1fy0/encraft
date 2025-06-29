package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/x/term"
	"github.com/electr1fy0/encraft/storage"
	"github.com/spf13/cobra"
)

var initCmd *cobra.Command = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new password vault",
	Long:  "Creates a new vault (encrypted) inside home directory",
	Run: func(cmd *cobra.Command, args []string) {
		exists, err := storage.VaultExists()
		if err != nil {
			fmt.Println("Error checking vault presence: ", err)
			os.Exit(1)
		}
		if exists {
			fmt.Println("Vault exists already. Add passwords")
			return
		}

		fmt.Print("Enter master password: ")
		passwordBytes, err := term.ReadPassword(0)
		if err != nil {
			fmt.Println("Error reading password: ", err)
			os.Exit(1)
		}

		password := string(passwordBytes)
		if len(password) < 8 {
			fmt.Println("Enter password with length > 8")
			os.Exit(1)
		}

		fmt.Println("Confirm master password: ")
		confirmBytes, err := term.ReadPassword(0)
		if err != nil {
			fmt.Println("Error reading confirmation password: ", err)
			os.Exit(1)
		}

		if password != string(confirmBytes) {
			fmt.Println("Passwords don't match.")
			os.Exit(1)
		}

		vault := storage.NewVault()

		err = storage.SaveVault(vault, password)

		if err != nil {
			fmt.Println("Error saving the vault", err)
			os.Exit(1)
		}

		fmt.Println("Vault initialized successfully")
		vaultPath, err := storage.GetVaultPath()
		if err != nil {
			fmt.Println("Error getting the vault path: ", err)
			os.Exit(1)
		}
		fmt.Println("Vault created at: ", vaultPath)

	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
