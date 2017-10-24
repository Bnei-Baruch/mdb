package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var batchCmd = &cobra.Command{
	Use:   "batch",
	Short: "run some batch command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use one of the subcommands")
	},
}

func init() {
	RootCmd.AddCommand(batchCmd)
}
