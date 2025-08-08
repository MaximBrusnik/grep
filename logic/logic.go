package logic

import (
	"bufio"
	"fmt"
	"grep/config"
	"os"
	"regexp"
	"strings"
)

func CreateMatcher(cfg *config.Config) func(string) bool {
	if cfg.Fixed {
		if cfg.IgnoreCase {
			pattern := strings.ToLower(cfg.Pattern)
			return func(line string) bool {
				return strings.Contains(strings.ToLower(line), pattern)
			}
		}
		return func(line string) bool {
			return strings.Contains(line, cfg.Pattern)
		}
	}

	pattern := cfg.Pattern
	if cfg.IgnoreCase {
		pattern = "(?i)" + pattern
	}
	re := regexp.MustCompile(pattern)
	return func(line string) bool {
		return re.MatchString(line)
	}
}

func ProcessFile(filename string, cfg *config.Config, matcher func(string) bool) (int, error) {
	var file *os.File
	if filename == "" {
		file = os.Stdin
	} else {
		var err error
		file, err = os.Open(filename)
		if err != nil {
			return 0, fmt.Errorf("error opening file: %v", err)
		}
		defer file.Close()
	}

	scanner := bufio.NewScanner(file)
	lineNum := 0
	count := 0
	var beforeLines []string
	afterLinesToPrint := 0

	printLine := func(line string, num int) {
		if cfg.LineNum {
			fmt.Printf("%d:%s\n", num, line)
		} else {
			fmt.Println(line)
		}
	}

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		matched := matcher(line)

		if (matched && !cfg.Invert) || (!matched && cfg.Invert) {
			for i, prevLine := range beforeLines {
				printLine(prevLine, lineNum-len(beforeLines)+i)
			}
			beforeLines = nil

			if !cfg.Count {
				printLine(line, lineNum)
			}
			count++

			afterLinesToPrint = cfg.After
		} else if afterLinesToPrint > 0 {
			if !cfg.Count {
				printLine(line, lineNum)
			}
			afterLinesToPrint--
		} else if cfg.Before > 0 {
			beforeLines = append(beforeLines, line)
			if len(beforeLines) > cfg.Before {
				beforeLines = beforeLines[1:]
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("error reading file: %v", err)
	}

	return count, nil
}
