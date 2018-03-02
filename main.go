package main

import (
	"fmt"
	"flag"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/xwb1989/sqlparser"

	"github.com/Chyroc/phpmyadmin-cli/internal"
)

var currentDB string

func execSQL(sql string) {
	sql = strings.TrimRight(sql, ";")
	sqlUpper := strings.ToUpper(sql)

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

		if sqlUpper == "SHOW DATABASES" || sqlUpper == "SHOW TABLES" {
			l := internal.ToList(selection)
			internal.FormatList(l[0], l[1:])
		}
	}
}

func executor(in string) {
	execSQL(in)
}

func completer(in prompt.Document) []prompt.Suggest {
	var suggest []prompt.Suggest
	for _, v := range internal.MySQLKeywords {
		suggest = append(suggest, prompt.Suggest{Text: v})
	}
	return prompt.FilterHasPrefix(suggest, in.GetWordBeforeCursor(), true)
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
