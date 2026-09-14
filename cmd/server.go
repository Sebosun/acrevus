package cmd

import (
	"sebosun/acrevus-go/server"

	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "starts the server",
	Run: func(cmd *cobra.Command, args []string) {
		server.StartServer()
	},
}
