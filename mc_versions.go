package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func checkMinecraftVersions() string {
	appdata := os.Getenv("APPDATA")
	versionsPath := filepath.Join(appdata, ".minecraft", "versions")

	entries, err := os.ReadDir(versionsPath)
	if err != nil {
		return "[11] ~ Couldt load folder: " + err.Error()
	}

	var sus []string
	allowedRegex := regexp.MustCompile(`^\d+\.\d+\.\d+$`)

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !allowedRegex.MatchString(name) {
			sus = append(sus, name)
		}
	}

	if len(sus) == 0 {
		return "[11] ~ Didnt find any suspicous version folder"
	}

	var builder strings.Builder
	builder.WriteString("[11] ~ Suspicous folders in .minecraft\\versions:\n")
	for _, s := range sus {
		builder.WriteString("      - " + s + "\n")
	}
	return builder.String()
}
