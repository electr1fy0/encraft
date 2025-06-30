package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/charmbracelet/x/term"
	"github.com/electr1fy0/encraft/storage"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [name]",
	Short: "Add password to the vault",
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

		// fmt.Print("Enter confirmation password: ")

		// confirmBytes, err := term.ReadPassword(0)
		// if err != nil {
		// 	fmt.Println("Error reading confirmation password")
		// 	os.Exit(1)
		// }

		// if password != string(confirmBytes) {
		// 	fmt.Println("Password don't match.")
		// 	os.Exit(1)
		// }

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
			fmt.Println("Entry already exists")
			os.Exit(1)
		}

		fmt.Print("Entry password: ")
		entryPassBytes, err := term.ReadPassword(0)
		if err != nil {
			fmt.Println("Couldn't read entry password")
			os.Exit(1)
		}

		entryPass := string(entryPassBytes)

		if entryPass == "" {
			fmt.Println("Password cannot be empty.")
			os.Exit(1)
		}

		entry := &storage.Entry{
			Name:     name,
			Password: entryPass,
		}

		vault.AddEntry(entry)
		err = storage.SaveVault(vault, password)
		if err != nil {
			fmt.Println("Error saving vault")
			os.Exit(1)
		}
		fmt.Println("Password added to vault")
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
