package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export some data from MDB",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use one of the subcommands")
	},
}

func init() {
	RootCmd.AddCommand(exportCmd)
}
