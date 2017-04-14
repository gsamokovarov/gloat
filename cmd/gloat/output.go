package main

import (
	"fmt"
	"os"
)

func Outf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stdout, format, a...)
}

func Exitf(code int, format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(code)
}
