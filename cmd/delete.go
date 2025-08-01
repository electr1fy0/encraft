package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/charmbracelet/x/term"
	"github.com/electr1fy0/encraft/storage"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "delete a password from the vault",
	Args:  cobra.ExactArgs(1),

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
		if err != nil {
			fmt.Println("Error reading password")
			os.Exit(1)

		}
		password := string(passwordBytes)

		vault, err := storage.LoadVault(password)
		if err != nil {
			fmt.Println("Error loading the vault")
			os.Exit(1)
		}

		reader := bufio.NewReader(os.Stdin)

		var name string
		if len(args) > 0 {
			name = args[0]
		} else {
			fmt.Print("Entry name: ")
			name, err = reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading entry name: ", err)
			}
		}

		if name == "" {
			println("Entry name cant be empty, mate")
			os.Exit(1)
		}

		if _, exists := vault.GetEntry(name); exists {
			vault.DeleteEntry(name)
		} else {
			fmt.Println("Entry does not exist")
			os.Exit(1)
		}

		err = storage.SaveVault(vault, password)
		if err != nil {
			fmt.Println("Error saving vault")
			os.Exit(1)
		}
		fmt.Println("Password removed from vault")
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
