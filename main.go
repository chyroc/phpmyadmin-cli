package main

import (
	"fmt"

	"github.com/c-bata/go-prompt"
	"github.com/xwb1989/sqlparser"
	"flag"
	"github.com/Chyroc/phpmyadmin-cli/internal"
	"strings"
)

var currentDB string

func execSQL(sql string) {
	sql = strings.TrimRight(sql, ";")

	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	switch  stmt.(type) {

	case *sqlparser.Use:
		dbs := strings.Split(sql, " ")
		currentDB = dbs[len(dbs)-1]
		fmt.Printf("Database changed: %s.\n", currentDB)
	default:
		selection, err := internal.Request(url, currentDB, sql)
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}

		if sql == "show databases" {
			internal.FormatList("Databases", internal.ToList(selection))
		} else if sql == "show tables"{
			internal.FormatList("Table", internal.ToList(selection))
		}

	}
}

func executor(in string) {
	execSQL(in)
}

func completer(in prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "use"},
		{Text: "insert"},
		{Text: "delete"},
		{Text: "select"},
		{Text: "update"},
		{Text: "databases"},
	}
	return prompt.FilterHasPrefix(s, in.GetWordBeforeCursor(), true)
}

var url string

func initConfig() {
	u := flag.String("url", "", "phpmyadmin url")
	flag.Parse()

	url = *u
}

func main() {
	initConfig()

	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix(">>> "),
		prompt.OptionTitle("phpmyadmin cli"),
	)
	p.Run()
}
