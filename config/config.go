package config

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	After      int
	Before     int
	Context    int
	Count      bool
	IgnoreCase bool
	Invert     bool
	Fixed      bool
	LineNum    bool
	Pattern    string
	Files      []string
}

func ParseFlags() *Config {
	cfg := &Config{}

	flag.IntVar(&cfg.After, "A", 0, "Print +N lines after each match")
	flag.IntVar(&cfg.Before, "B", 0, "Print +N lines before each match")
	flag.IntVar(&cfg.Context, "C", 0, "Print ±N lines around each match")
	flag.BoolVar(&cfg.Count, "c", false, "Print only a count of selected lines")
	flag.BoolVar(&cfg.IgnoreCase, "i", false, "Ignore case distinctions")
	flag.BoolVar(&cfg.Invert, "v", false, "Selected lines are those not matching any of the specified patterns")
	flag.BoolVar(&cfg.Fixed, "F", false, "Interpret patterns as fixed strings, not regular expressions")
	flag.BoolVar(&cfg.LineNum, "n", false, "Prefix each line of output with the 1-based line number within its input file")

	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Usage: grep [OPTIONS] PATTERN [FILE...]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	cfg.Pattern = args[0]
	if len(args) > 1 {
		cfg.Files = args[1:]
	}

	if cfg.Context > 0 {
		cfg.After = cfg.Context
		cfg.Before = cfg.Context
	}

	return cfg
}
