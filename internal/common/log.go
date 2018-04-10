package common

import (
	"github.com/fatih/color"
)

func Info(fotmat string, a ...interface{}) {
	color.Green(fotmat, a...)
}

func Warn(fotmat string, a ...interface{}) {
	color.Red(fotmat, a...)
}

func Error(err error) {
	color.Red("%s\n", err)
}
