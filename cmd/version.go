package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/Bnei-Baruch/mdb/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of MDB",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("BB archive Metadata Database version %s\n", version.Version)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
