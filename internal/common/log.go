package common

import (
	"github.com/fatih/color"
)

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
