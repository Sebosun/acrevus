package cmd

import (
	"fmt"
	"net/url"
	"os"
	"sebosun/acrevus-go/analyzer"
	"sebosun/acrevus-go/fetcher"
	"sebosun/acrevus-go/storage"

	"github.com/spf13/cobra"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch [url]",
	Short: "Fetch a URL",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			fmt.Println("Multiple url links included, not supported yet")
			os.Exit(0)
		}
		fmt.Println("Fetch start...")
		link := args[0]
		_, err := url.ParseRequestURI(link)
		if err != nil {
			fmt.Println("Invalid url provided")
			os.Exit(1)
		}

		isSaved, err := storage.IsURLSaved(link)
		if err != nil {
			fmt.Printf("error reading from articles json. %i \n", err)
			os.Exit(1)
		}
		if isSaved {
			fmt.Println("Article already saved")
			os.Exit(0)
		}

		fetcher.InitFetcher(link)
	},
}

// For now it's mostly internal for running command out of cli
var parseURL = &cobra.Command{
	Use:     "analyze [url]",
	Aliases: []string{"a"},
	Short:   "Analyze a URL",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running analyzer")
		if len(args) > 1 {
			fmt.Println("Multiple url links included, not supported yet")
			os.Exit(0)
		}
		link := args[0]
		_, err := url.ParseRequestURI(link)
		if err != nil {
			fmt.Println("Invalid url provided")
			os.Exit(1)
		}

		err = analyzer.Run(link)
		if err != nil {
			fmt.Println("Error running analyzer")
			os.Exit(1)
		}
	},
}
