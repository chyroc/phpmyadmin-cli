package common

import (
	"io"
	"log"
	"os"

	"github.com/fatih/color"
)

var f io.Writer

func InitLog(logPath string) {
	var err error
	f, err = os.OpenFile(logPath, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
}

func Debug(format string, a ...interface{}) {
	if IsDebug1 {
		color.HiBlack(format, a...)
	}
}

func Info(format string, a ...interface{}) {
	color.Green(format, a...)
}

func Warn(format string, a ...interface{}) {
	color.Red(format, a...)
}

func Error(err error) {
	color.Red("%s\n", err)
}

func Exit(err error) {
	Error(err)
	os.Exit(1)
}

func Logf(format string, a ...interface{}) {
	if f == nil {
		f = os.Stdout
	}
	l := log.New(f, "", log.LstdFlags)
	l.Printf(format, a...)
}
