package main

import (
	"os"
	"strings"
)

var suspiciousKeywords = []string{
	"vape", "autoclicker", "clicker", "cheat", "raven", "injector", "ghost",
}

func checkPrefetch() string {
	prefetchPath := `C:\Windows\Prefetch`
	entries, err := os.ReadDir(prefetchPath)
	if err != nil {
		return "[3] ~ Prefetch folder not accessible"
	}

	var found []string

	for _, entry := range entries {
		name := strings.ToLower(entry.Name())
		for _, keyword := range suspiciousKeywords {
			if strings.Contains(name, keyword) {
				found = append(found, entry.Name())
				break
			}
		}
	}

	if len(found) == 0 {
		return "[3] ~ No suspicious prefetch files found"
	}

	var builder strings.Builder
	builder.WriteString("[3] ~ Suspicious prefetch files:\n")
	for _, match := range found {
		builder.WriteString("      - " + match + "\n")
	}

	return builder.String()
}
