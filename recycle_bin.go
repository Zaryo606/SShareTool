package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"time"
)

func checkRecycleBin() string {
	basePath := `C:\$Recycle.Bin`
	entries, err := os.ReadDir(basePath)
	if err != nil {
		return "[2] ~ Recycle Bin folder not found"
	}

	var latestMod time.Time
	var usedFolder string

	for _, entry := range entries {
		fullPath := filepath.Join(basePath, entry.Name())
		info, err := os.Stat(fullPath)
		if err != nil {
			continue
		}
		if info.ModTime().After(latestMod) {
			latestMod = info.ModTime()
			usedFolder = fullPath
		}
	}

	if usedFolder == "" {
		return "[2] ~ No valid Recycle Bin folders found"
	}

	currentUser, userErr := user.Current()
	username := "unknown"
	if userErr == nil {
		username = currentUser.Username
	}

	return fmt.Sprintf(
		"[2] ~ Recycle Bin info:\n"+
			"      - Last modified: %s\n"+
			"      - Path: %s\n"+
			"      - Current user: %s",
		latestMod.Format("02.01.2006 15:04"),
		usedFolder,
		username,
	)
}
