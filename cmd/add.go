package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/x/term"
	"github.com/electr1fy0/encraft/storage"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [name]",
	Short: "Add password to the vault",
	Args:  cobra.MaximumNArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		exists, err := storage.VaultExists()
		if err != nil {
			fmt.Println("Error checking vault's presence:", err)
			os.Exit(1)
		}

		if !exists {
			fmt.Println("No vault found. Run `encraft init` first.")
			os.Exit(1)
		}

		fmt.Print("Enter master password: ")
		passwordBytes, err := term.ReadPassword(0)
		if err != nil {
			fmt.Println("\nError reading master password")
			os.Exit(1)
		}
		fmt.Println()
		password := string(passwordBytes)

		vault, err := storage.LoadVault(password)
		if err != nil {
			fmt.Println("Error loading vault:", err)
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
				fmt.Println("Error reading entry name:", err)
				os.Exit(1)
			}
		}

		name = strings.TrimSpace(name)

		if name == "" {
			fmt.Println("Entry name cannot be empty.")
			os.Exit(1)
		}

		if _, exists := vault.GetEntry(name); exists {
			fmt.Println("Entry with this name already exists.")
			os.Exit(1)
		}

		fmt.Print("Entry password: ")
		entryPassBytes, err := term.ReadPassword(0)
		if err != nil {
			fmt.Println("\nCouldn't read entry password")
			os.Exit(1)
		}
		fmt.Println()
		entryPass := strings.TrimSpace(string(entryPassBytes))

		if entryPass == "" {
			fmt.Println("Password cannot be empty.")
			os.Exit(1)
		}

		fmt.Print("URL (optional): ")
		url, _ := reader.ReadString('\n')
		url = strings.TrimSpace(url)

		fmt.Print("Notes (optional): ")
		notes, _ := reader.ReadString('\n')
		notes = strings.TrimSpace(notes)

		entry := &storage.Entry{
			Name:     name,
			Password: entryPass,
			URL:      url,
			Notes:    notes,
		}

		vault.AddEntry(entry)

		err = storage.SaveVault(vault, password)
		if err != nil {
			fmt.Println("Error saving vault:", err)
			os.Exit(1)
		}

		fmt.Println("✅ Entry added successfully:")
		fmt.Println("  Name  :", name)
		if url != "" {
			fmt.Println("  URL   :", url)
		}
		if notes != "" {
			fmt.Println("  Notes :", notes)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
