package cmd

import (
	"sebosun/acrevus-go/repl"

	"github.com/spf13/cobra"
)

var uiCmd = &cobra.Command{
	Use:   "ui",
	Short: "Show Acrevus UI",
	Run: func(cmd *cobra.Command, args []string) {
		repl.InitTea()
		// Your root command logic here
	},
}
