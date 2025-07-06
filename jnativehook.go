package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func checkTempJNativeHook() string {
	tempDir := os.TempDir()
	var matches []string

	err := filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		if strings.Contains(strings.ToLower(info.Name()), "jnativehook") {
			modTime := info.ModTime().Format("02.01.2006 15:04:05")
			matches = append(matches, fmt.Sprintf("      - %s\n        Last modified: %s", path, modTime))
		}
		return nil
	})

	if err != nil {
		return fmt.Sprintf("[6] ~ Error scanning %%TEMP%%: %v", err)
	}

	if len(matches) == 0 {
		return "[6] ~ No 'jnativehook' files found in %TEMP%"
	}

	builder := strings.Builder{}
	builder.WriteString("[6] ~ Suspicious 'jnativehook' files in %TEMP%:\n")
	for _, m := range matches {
		builder.WriteString(m + "\n")
	}
	return builder.String()
}
