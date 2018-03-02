package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/xwb1989/sqlparser"

	"github.com/Chyroc/phpmyadmin-cli/internal"
)

var currentDB string
var url string
var historyPath string
var history []string

func addHistory(word string) {
	f, err := os.OpenFile(historyPath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return
	}

	defer f.Close()
	f.WriteString("\n" + word)
}

func execSQL(sql string) {
	sql = strings.TrimRight(sql, ";")
	addHistory(sql)

	sqlUpper := strings.ToUpper(sql)

	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	selection, err := internal.Request(url, currentDB, sql)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	switch stmt.(type) {

	case *sqlparser.Show:
		if sqlUpper == "SHOW TABLES" && currentDB == "" {
			fmt.Printf(`(1046, u'No database selected')` + "\n")
			return
		}

		l := internal.ToList(selection)
		internal.FormatList(l[0], l[1:])
		return
	case *sqlparser.Use:
		dbs := strings.Split(sql, " ")
		db := dbs[len(dbs)-1]
		if selection == nil {
			fmt.Printf(`(1049, u"Unknown database '%s'")`+"\n", db)
			return
		}

		currentDB = db

		LivePrefixState.LivePrefix = db + " >>> "
		LivePrefixState.IsEnable = true

		fmt.Printf("Database changed: %s.\n", currentDB)
	default:
		fmt.Printf("select: %s\n", selection.Text())
	}
}

var LivePrefixState struct {
	LivePrefix string
	IsEnable   bool
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

func initConfig() {
	u := flag.String("url", "", "phpmyadmin url")
	hPath := flag.String("h", ".phpmyadmin_cli_history", "phpmyadmin url")
	flag.Parse()

	url = *u
	historyPath = *hPath

	body, err := ioutil.ReadFile(historyPath)
	if err != nil {
		panic(err)
	}
	history = strings.Split(string(body), "\n")
}

func main() {
	initConfig()

	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix(">>> "),
		prompt.OptionTitle("phpmyadmin cli"),
		prompt.OptionHistory(history),
		prompt.OptionLivePrefix(func() (prefix string, useLivePrefix bool) {
			return LivePrefixState.LivePrefix, LivePrefixState.IsEnable
		}),
	)
	p.Run()
}
