package main

import (
	"fmt"
	"grep/config"
	"grep/logic"
	"os"
)

func main() {
	cfg := config.ParseFlags()
	matcher := logic.CreateMatcher(cfg)

	totalCount := 0

	if len(cfg.Files) == 0 {
		count, err := logic.ProcessFile("", cfg, matcher)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		totalCount += count
	} else {
		for _, file := range cfg.Files {
			count, err := logic.ProcessFile(file, cfg, matcher)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %v\n", file, err)
				continue
			}
			totalCount += count
		}
	}

	if cfg.Count {
		fmt.Println(totalCount)
	}
}
