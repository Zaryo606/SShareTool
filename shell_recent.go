package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func checkRecentMacros() string {
	recentPath := filepath.Join(os.Getenv("APPDATA"), "Microsoft", "Windows", "Recent")
	suspiciousWords := []string{"macro", "autoclick", "clicker", "bot", "jnativehook", "script"}
	var suspiciousMatches []string
	var usbMatches []string

	var fileCount int
	var oldest, newest time.Time

	err := filepath.Walk(recentPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		fileCount++
		modTime := info.ModTime()
		if oldest.IsZero() || modTime.Before(oldest) {
			oldest = modTime
		}
		if newest.IsZero() || modTime.After(newest) {
			newest = modTime
		}

		name := info.Name()
		lowerName := strings.ToLower(name)

		// Pokud je název podezřelý
		for _, word := range suspiciousWords {
			if strings.Contains(lowerName, word) && strings.HasSuffix(lowerName, ".lnk") {
				suspiciousMatches = append(suspiciousMatches,
					fmt.Sprintf("      - %s\n        Last modified: %s", name, modTime.Format("02.01.2006 15:04:05")))
				return nil
			}
		}

		// Pokud jde o zástupce disku (např. "VIT KOLATOR (D).lnk")
		if strings.HasSuffix(lowerName, ".lnk") && looksLikeDriveShortcut(name) {
			usbMatches = append(usbMatches,
				fmt.Sprintf("      - %s\n        Last modified: %s", name, modTime.Format("02.01.2006 15:04:05")))
		}

		return nil
	})

	if err != nil {
		return fmt.Sprintf("[7] ~ Error scanning recent folder: %v", err)
	}

	builder := strings.Builder{}
	builder.WriteString("[7] ~ Recent folder scan results:\n")

	// PODEZŘELÉ ZÁSTUPCE
	if len(suspiciousMatches) > 0 {
		builder.WriteString("    → Suspicious files found in 'Recent':\n")
		for _, m := range suspiciousMatches {
			builder.WriteString(m + "\n")
		}
	} else {
		builder.WriteString("    → No suspicious macro/autoclicker references found in recent files\n")
	}

	// USB ZÁSTUPCE
	if len(usbMatches) > 0 {
		builder.WriteString("    → Shortcuts that indicate USB drive was opened:\n")
		for _, m := range usbMatches {
			builder.WriteString(m + "\n")
		}
	} else {
		builder.WriteString("    → No shortcuts to external drives detected\n")
	}

	return builder.String()
}

// Detects names like "VIT KOLATOR (D)" or "Backup (E)"
func looksLikeDriveShortcut(name string) bool {
	name = strings.ToUpper(name)
	for letter := 'D'; letter <= 'Z'; letter++ {
		if strings.Contains(name, fmt.Sprintf("(%c)", letter)) {
			return true
		}
	}
	return false
}
