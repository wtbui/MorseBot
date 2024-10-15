package options

import (
	flag "github.com/spf13/pflag"
)

// Bot Deployment Options
type Options struct {
	Version bool
	Help bool
	Verbose bool
}

var mb_options = &Options{}

func addFlags() {
	flag.BoolVar(&mb_options.Version, "version", false, "Version of Morse Bot")
	flag.BoolVar(&mb_options.Verbose, "debug", false, "Start in debug mode")
}

func ParseOptions(args []string) (*Options, error) {
	addFlags()
	flag.Parse()

	return mb_options, nil
}
