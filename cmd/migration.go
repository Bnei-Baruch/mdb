package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"time"
)

func init() {
	command := &cobra.Command{
		Use:   "migration <name>",
		Short: "Create migration file",
		Long:  "Create new migration file",
		Run:   migrationFn,
	}

	RootCmd.AddCommand(command)
}

func migrationFn(cmd *cobra.Command, args []string) {
	const template = `
-- MDB generated migration file
-- rambler up

-- rambler down

`
	var migrationName string

	if len(args) > 0 {
		t := time.Now()
		timestamp := t.Format("2006-01-02_12-04-05_")
		migrationName = "./migrations/" + timestamp + args[0] + ".sql"
	} else {
		fmt.Println("Please specify migration name")
		return
	}
	if err := ioutil.WriteFile(migrationName, []byte(template), 0644); err != nil {
		panic(err)
	} else {
		fmt.Printf("Migration file %s was created successfuly\n", migrationName)
	}
}
