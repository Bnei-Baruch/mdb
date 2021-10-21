package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/volatiletech/sqlboiler/boil"

	"github.com/Bnei-Baruch/mdb/batch"
	"github.com/Bnei-Baruch/mdb/utils"
)

var regexpReplacerCmd = &cobra.Command{
	Use:   "reg_replace",
	Short: "replace string with regular expression by table name + column name",
	Run:   regexpReplacerFn,
}

func init() {
	RootCmd.AddCommand(regexpReplacerCmd)
}

func regexpReplacerFn(cmd *cobra.Command, args []string) {
	//log.SetLevel(log.DebugLevel)
	boil.DebugMode = true
	if len(args) < 10 {
		fmt.Print(`You need enter 4 arguments:\n 
			1 - table name,\n
			2 - column name,\n
			3 - regular expression,\n 
			4 - string for replace\n
			For example: ./mdb reg_replace "blog_posts" "content" "(http://.{0,5}youtube)" "https://www.youtube"
		`)
		return
	}

	replacer := batch.RegexpReplacer{
		TableName: args[0],
		ColName:   args[1],
		RegStr:    args[2],
		NewStr:    args[3],
	}
	utils.Must(replacer.Init())
	replacer.Do()
}
