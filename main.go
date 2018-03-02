package main

import (
	"fmt"

	"github.com/c-bata/go-prompt"
	"github.com/xwb1989/sqlparser"
)

func execSQL(sql string) {
	stmt, err := sqlparser.Parse(sql)
	if err != nil {
	}

	switch stmt := stmt.(type) {
	case *sqlparser.Select:
		_ = stmt
	case *sqlparser.Insert:
	}
}

func executor(in string) {
	fmt.Println("Your input: " + in)
	execSQL(in)
}

func completer(in prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "use"},
		{Text: "insert"},
		{Text: "delete"},
		{Text: "select"},
		{Text: "update"},
	}
	return prompt.FilterHasPrefix(s, in.GetWordBeforeCursor(), true)
}

func main() {
	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix(">>> "),
		prompt.OptionTitle("phpmyadmin cli"),
	)
	p.Run()
}
