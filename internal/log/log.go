package log

import (
	"fmt"
	"os"
)

const (
	Debug = "DEBUG"
	Info  = "INFO"
	Warn  = "WARN"
	Error = "ERROR"
)

const (
	reset  = "\033[0m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	blue   = "\033[34m"
)

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

	content := fmt.Sprintf(format, args...)
	fmt.Fprintf(os.Stdout, "\x1b[1m%s[%s]%s %s\n", colour, level, reset, content)
}
