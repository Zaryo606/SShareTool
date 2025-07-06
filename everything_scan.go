package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func checkSuspiciousNamesOnCDrive() string {
	suspiciousWords := []string{"macro", "autoclick", "clicker", "injecter.exe", "jnativehook"}
	excludedDirs := map[string]bool{
		`C:\Windows`:                   true,
		`C:\ProgramData`:               true,
		`C:\$Recycle.Bin`:              true,
		`C:\System Volume Information`: true,
		`C:\Recovery`:                  true,
		`C:\PerfLogs`:                  true,
	}

	var results []string

	err := filepath.Walk(`C:\`, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			if excludedDirs[path] {
				return filepath.SkipDir
			}
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
		return fmt.Sprintf("[8] ~ Error scanning C:\\ drive: %v", err)
	}

	if len(results) == 0 {
		return "[8] ~ No suspicious file/folder names found on C:\\"
	}

	var builder strings.Builder
	builder.WriteString("[8] ~ Suspicious file/folder names found on C:\\\n")
	for _, line := range results {
		builder.WriteString(line + "\n")
	}
	return builder.String()
}
