package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func checkRecentMacros() string {
	recentPath := filepath.Join(os.Getenv("APPDATA"), "Microsoft", "Windows", "Recent")
	suspiciousWords := []string{"macro", "autoclick", "clicker", "bot", "jnativehook", "script"}
	var matches []string

	err := filepath.Walk(recentPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		lowerName := strings.ToLower(info.Name())
		for _, word := range suspiciousWords {
			if strings.Contains(lowerName, word) {
				modTime := info.ModTime().Format("02.01.2006 15:04:05")
				matches = append(matches, fmt.Sprintf("      - %s\n        Last modified: %s", info.Name(), modTime))
				break
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Sprintf("[7] ~ Error scanning recent folder: %v", err)
	}

	if len(matches) == 0 {
		return "[7] ~ No suspicious macro/autoclicker references found in recent files"
	}

	builder := strings.Builder{}
	builder.WriteString("[7] ~ Suspicious files found in 'Recent':\n")
	for _, m := range matches {
		builder.WriteString(m + "\n")
	}
	return builder.String()
}
