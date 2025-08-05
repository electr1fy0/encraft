package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/charmbracelet/x/term"
	"github.com/electr1fy0/encraft/storage"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "Show all entries",
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, args []string) {
		exists, err := storage.VaultExists()
		if err != nil {
			fmt.Println("Error checking vault presence:", err)
			os.Exit(1)
		}
		if !exists {
			fmt.Println("No vault found. Initialize first.")
			os.Exit(1)
		}

		fmt.Print("Enter master password: ")
		passwordBytes, err := term.ReadPassword(0)
		if err != nil {
			fmt.Println("Error reading password")
			os.Exit(1)
		}
		fmt.Println() // for newline after password input

		password := string(passwordBytes)
		vault, err := storage.LoadVault(password)
		if err != nil {
			fmt.Println("Error loading vault:", err)
			os.Exit(1)
		}

		names := vault.ListEntries()
		sort.Strings(names)

		if len(names) == 0 {
			fmt.Println("No entries in the vault.")
			return
		}

		data := [][]string{}
		for i, name := range names {
			entry, _ := vault.GetEntry(name)

			created := entry.CreatedAt.Format("2006-01-02 15:04")
			updated := entry.UpdatedAt.Format("2006-01-02 15:04")
			notes := entry.Notes
			if len(notes) > 20 {
				notes = notes[:20] + "..."
			}

			data = append(data, []string{
				fmt.Sprintf("%d", i+1),
				entry.Name,
				entry.URL,
				created,
				updated,
				notes,
			})
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.Header([]string{"#", "Name", "URL", "Created", "Updated", "Notes"})
		table.Bulk(data)
		table.Render()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
