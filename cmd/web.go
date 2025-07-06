package cmd

import (
	"fmt"
	"os"

	"github.com/electr1fy0/encraft/server"
	"github.com/spf13/cobra"
)

var (
	port string
	host string
)

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Launch the web interface",
	Run: func(cmd *cobra.Command, args []string) {
		// exists, err := storage.VaultExists()
		// if err != nil {
		// 	fmt.Println("Error checking vault: ", err)
		// 	os.Exit(1)
		// }

		// if !exists {
		// 	fmt.Println("No vault found. Run 'secrets init' first.")
		// 	os.Exit(1)
		// }
		fmt.Printf("Starting web server on:%s:%s", host, port)
		println()

		server := server.NewServer()
		if err := server.Start(host + ":" + port); err != nil {
			fmt.Println("Error starting the server", err)
			os.Exit(1)
		}

	},
}

func init() {
	webCmd.Flags().StringVarP(&port, "port", "p", "8080", "port to run the web server on")
	webCmd.Flags().StringVarP(&host, "host", "l", "localhost", "host to bind to")
	rootCmd.AddCommand(webCmd)
}
