package log

import (
	"fmt"
	"os"
)

const (
	Debug = "D"
	Info  = "I"
	Warn  = "W"
	Error = "E"
)

const (
	reset  = "\033[0m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	blue   = "\033[34m"
	purple = "\033[35m"
	orange = "\033[36m"
	grey   = "\033[37m"
)

var Colours = []string{purple, orange, grey}

func Log(lvl string, format string, args ...any) {
	var col string

	switch lvl {
	case Debug:
		col = blue
	case Info:
		col = green
	case Warn:
		col = yellow
	case Error:
		col = red
	default:
		col = reset
	}

	txt := fmt.Sprintf(format, args...)
	fmt.Fprintf(os.Stdout, "\x1b[1m%s%s.%s %s\n", col, lvl, reset, txt)
}

func Abort(fmt string, args ...any) {
	Log(Error, fmt, args...)
	os.Exit(1)
}

func PrettyContent(cmds []string, conds []string) {
	for idx, command := range cmds {
		Log(
			Info,
			"%s = \x1b[1m%s\033[0m %d. \x1b[1m%s\033[0m",
			conds[0],
			conds[1],
			idx,
			command,
		)
	}
}
