// Package cmd handles inline args
package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{Use: "app"}

func Execute() {
	rootCmd.AddCommand(fetchCmd)
	rootCmd.AddCommand(uiCmd)
	rootCmd.Execute()
}
