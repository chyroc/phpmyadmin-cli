package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/xwb1989/sqlparser"

	"github.com/Chyroc/phpmyadmin-cli/internal/common"
	"github.com/Chyroc/phpmyadmin-cli/internal/phpmyadmin"
	"github.com/Chyroc/phpmyadmin-cli/internal/show"
)

var currentDB string
var url string
var historyPath string
var history []string
var help bool
var prune bool
var list bool
var server string
var ErrNotSetServer = errors.New("not set server")

func getHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	return usr.HomeDir
}

func addHistory(word string) {
	f, err := os.OpenFile(historyPath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return
	}

	defer f.Close()
	f.WriteString("\n" + word)
}

func setServer(sql string) error {
	sqls := strings.Split(strings.TrimSpace(strings.ToLower(sql)), " ")
	if len(sqls) == 2 && sqls[0] == "set" {
		s, err := phpmyadmin.DefaultPhpmyadmin.GetServerList(url)
		if err != nil {
			return err
		}

		for _, v := range s.S {
			if v.ID == sqls[1] {
				server = v.ID
				setPrefix()
				common.Info("server[%s] seted", v.ID)
				return nil
			}
		}
		err = fmt.Errorf("server seted[%s] is invalid", sqls[1])
		return err
	}

	if server == "" {
		err := fmt.Errorf("请选择一个server; 输入`show servers`获取所有server; 输入`set <id>`设置server")
		return err
	}

	return ErrNotSetServer
}

func setPrefix() {
	promptPrefix := ""
	if server != "" && currentDB != "" {
		promptPrefix = "[" + server + "] " + currentDB + " >>> "
	} else if currentDB == "" {
		promptPrefix = "[" + server + "]" + " >>> "
	}

	LivePrefixState.LivePrefix = promptPrefix
	LivePrefixState.IsEnable = true
}

func execSQL(sql string) {
	if strings.ToLower(sql) == "show servers" {
		s, err := phpmyadmin.DefaultPhpmyadmin.GetServerList(url)
		if err != nil {
			common.Error(err)
		}
		s.Print()
		return
	}

	err := setServer(sql)
	if err != ErrNotSetServer {
		if err != nil {
			common.Error(err)
		}
		return
	}

	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		common.Warn("syntax error: %s %s\n", sql, err)
		return
	}

	result, err := phpmyadmin.DefaultPhpmyadmin.ExecSQL(server, currentDB, "", sql)
	if err != nil {
		common.Error(err)
		return
	}

	switch stmt.(type) {
	case *sqlparser.Use:
		dbs := strings.Split(sql, " ")
		db := dbs[len(dbs)-1]
		if result == nil {
			common.Warn(`(1049, u"Unknown database '%s'")`+"\n", db)
			return
		}

		currentDB = db
		setPrefix()

		common.Info("Database changed: %s.\n", currentDB)
		return
	default:
		if strings.ToUpper(sql) != "SHOW DATABASES" && currentDB == "" {
			common.Warn("(1046, u'No database selected')\n")
			return
		}

		html := string(result)
		show.ParseFromHTML(fmt.Sprintf("<table>%s</table>", html))
	}
}

var LivePrefixState struct {
	LivePrefix string
	IsEnable   bool
}

func executor(in string) {
	addHistory(in)
	sqls := strings.Split(in, ";")
	for _, s := range sqls {
		if s != "" {
			execSQL(s)
		}
	}
}

func completer(in prompt.Document) []prompt.Suggest {
	var suggest []prompt.Suggest
	for _, v := range common.MySQLKeywords {
		suggest = append(suggest, prompt.Suggest{Text: v})
	}

	suggest = append(suggest, prompt.Suggest{Text: "servers"})

	return prompt.FilterHasPrefix(suggest, in.GetWordBeforeCursor(), true)
}

func initConfig() {
	flag.StringVar(&url, "url", "", "phpmyadmin url")
	flag.StringVar(&historyPath, "history", getHomeDir()+"/.phpmyadmin_cli_history", "phpmyadmin url")
	flag.BoolVar(&help, "h", false, "show help")
	flag.BoolVar(&prune, "prune", false, "清理命令记录")
	flag.BoolVar(&list, "list", false, "获取server列表")
	flag.StringVar(&server, "server", "", "选择server")
	flag.Parse()

	phpmyadmin.DefaultPhpmyadmin.SetURI(url)

	if len(flag.Args()) > 0 {
		fmt.Printf("", flag.Args())
		// execSQL("use " + flag.Args()[0])
	}

	body, err := ioutil.ReadFile(historyPath)
	if err != nil {
		panic(err)
	}
	history = strings.Split(string(body), "\n")
}

func reverseStrings(input []string) []string {
	if len(input) == 0 {
		return input
	}
	return append(reverseStrings(input[1:]), input[0])
}

func shortHistory() error {
	body, err := ioutil.ReadFile(historyPath)
	if err != nil {
		return err
	}

	lines := reverseStrings(strings.Split(string(body), "\n"))

	var newLines []string
	var lineSet = make(map[string]bool)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !lineSet[line] {
			newLines = append(newLines, line)
			lineSet[line] = true
		}
	}

	return ioutil.WriteFile(historyPath, []byte(strings.Join(reverseStrings(newLines), "\n")), 0644)
}

func main() {
	initConfig()

	if help {
		fmt.Printf(`NAME:
   phpmyadmin-cli - access phpmyadmin from shell cli

USAGE:
   phpmyadmin-cli [global options] [arguments...]

GLOBAL OPTIONS:
   --url value      phpmyadmin url
   --prune          清理命令记录
   --server         选择server
   --list           获取server列表
   --history value  command history file (default: "%s")
   --help, -h       show help`+"\n", historyPath)
		return
	} else if prune {
		err := shortHistory()
		if err != nil {
			common.Error(err)
		}
		return
	}

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
