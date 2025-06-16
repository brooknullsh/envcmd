package log

import (
	"fmt"
	"os"
)

const (
	blue    = "\033[34m"
	magenta = "\033[35m"
	cyan    = "\033[36m"
	white   = "\033[37m"
)

var Colours = []string{blue, magenta, cyan, white}

func Info(format string, args ...any) {
	fmt.Printf("\x1b[1m\033[32mINFO\033[0m %s\n", fmt.Sprintf(format, args...))
}

func Warn(format string, args ...any) {
	fmt.Printf("\x1b[1m\033[33mWARN\033[0m %s\n", fmt.Sprintf(format, args...))
}

func Abort(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "\x1b[1m\033[31mERROR\033[0m %s\n", fmt.Sprintf(format, args...))
	os.Exit(1)
}
