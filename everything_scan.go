package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func checkSuspiciousNamesOnCDrive() string {
	suspiciousWords := []string{"macro", "autoclick", "clicker", "injecter.exe", "jnativehook"}
	var results []string

	for letter := 'C'; letter <= 'Z'; letter++ {
		drive := fmt.Sprintf("%c:\\", letter)
		if _, err := os.Stat(drive); os.IsNotExist(err) {
			continue
		}

		err := filepath.Walk(drive, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}

			name := strings.ToLower(info.Name())
			for _, word := range suspiciousWords {
				if strings.Contains(name, word) {
					results = append(results, "      - "+path)
					break
				}
			}
			return nil
		})

		if err != nil {
			results = append(results, fmt.Sprintf("      - Error scanning %s: %v", drive, err))
		}
	}

	if len(results) == 0 {
		return "[8] ~ No suspicious file/folder names found on any drive"
	}

	var builder strings.Builder
	builder.WriteString("[8] ~ Suspicious file/folder names found:\n")
	for _, line := range results {
		builder.WriteString(line + "\n")
	}
	return builder.String()
}
