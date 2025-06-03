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

func Log(level string, format string, args ...any) {
	var colour string

	switch level {
	case Debug:
		colour = blue
	case Info:
		colour = green
	case Warn:
		colour = yellow
	case Error:
		colour = red
	default:
		colour = reset
	}

	output := fmt.Sprintf(format, args...)
	fmt.Fprintf(os.Stdout, "\x1b[1m%s%s.%s %s\n", colour, level, reset, output)
}

func Abort(format string, args ...any) {
	Log(Error, format, args...)
	os.Exit(1)
}
