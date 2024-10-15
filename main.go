package main

import (
	"fmt"
	"os"

	mb "https://github.com/wtbui/MorseBot/cmd"
	options "https://github.com/wtbui/MorseBot/pkg/options"
)

var version = "0.1"

func exit(code int, err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	os.Exit(code)
}

func main() {
	opts, err := options.ParseFlags(os.Args[1:])
	if err != nil {
		exit(1, err)
		return
	}

	if opts.Version {
		fmt.Println(version)
		return
	}

	code, err := mb.Start(opts)
}
