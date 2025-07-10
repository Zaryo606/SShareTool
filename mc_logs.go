package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/process"
)

func checkMinecraftLogs() string {
	appdata := os.Getenv("APPDATA")
	logsDir := filepath.Join(appdata, ".minecraft", "logs")
	latestLog := filepath.Join(logsDir, "latest.log")

	javaProc, err := findJavawProcess()
	if err != nil {
		return "[10] ~ javaw.exe is not running"
	}
	startTime, err := javaProc.CreateTime()
	if err != nil {
		return "[10] ~ Failed to determine javaw.exe start time"
	}
	startTimeT := time.Unix(0, startTime*int64(time.Millisecond))

	info, err := os.Stat(latestLog)
	if err != nil {
		return "[10] ~ latest.log does not exist"
	}
	modTime := info.ModTime()

	if modTime.Before(startTimeT) {
		return "[10] ~ latest.log was not modified after Minecraft started → logging error"
	}

	activeUser, err := getLatestUserFromLog(latestLog)
	if err != nil {
		activeUser = "?? (not found)"
	}

	var builder strings.Builder
	builder.WriteString("[10] ~ Minecraft log analysis:\n")
	builder.WriteString("      - Current user: " + activeUser + "\n")

	users, err := extractUsersFromGzLogs(logsDir)
	if err == nil && len(users) > 0 {
		builder.WriteString("      - Other user accounts found in archived logs:\n")
		for user, files := range users {
			if user == activeUser {
				continue
			}
			builder.WriteString(fmt.Sprintf("         → %s (in logs: %s)\n", user, strings.Join(files, ", ")))
		}
	}

	return builder.String()
}

func findJavawProcess() (*process.Process, error) {
	processes, err := process.Processes()
	if err != nil {
		return nil, err
	}
	for _, p := range processes {
		name, err := p.Name()
		if err == nil && strings.EqualFold(name, "javaw.exe") {
			return p, nil
		}
	}
	return nil, fmt.Errorf("javaw.exe not found")
}

func getLatestUserFromLog(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var lastUser string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Successfully refreshed token for ") {
			start := strings.Index(line, "Successfully refreshed token for ") + len("Successfully refreshed token for ")
			if start < len(line) {
				lastUser = strings.TrimSpace(line[start:])
			}
		}
	}

	if lastUser != "" {
		return lastUser, nil
	}

	file.Seek(0, 0)
	scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Setting user: ") {
			start := strings.Index(line, "Setting user: ") + len("Setting user: ")
			if start < len(line) {
				lastUser = strings.TrimSpace(line[start:])
			}
		}
	}

	if lastUser != "" {
		return lastUser, nil
	}
	return "", fmt.Errorf("no user found in log")
}

func extractUsersFromGzLogs(folder string) (map[string][]string, error) {
	users := make(map[string][]string)

	entries, err := os.ReadDir(folder)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".gz") {
			continue
		}

		path := filepath.Join(folder, entry.Name())
		file, err := os.Open(path)
		if err != nil {
			continue
		}
		gzr, err := gzip.NewReader(file)
		if err != nil {
			file.Close()
			continue
		}
		scanner := bufio.NewScanner(gzr)
		for scanner.Scan() {
			line := scanner.Text()

			var prefix string
			if strings.Contains(line, "Successfully refreshed token for ") {
				prefix = "Successfully refreshed token for "
			} else if strings.Contains(line, "Setting user: ") {
				prefix = "Setting user: "
			} else {
				continue
			}

			start := strings.Index(line, prefix)
			if start == -1 {
				continue
			}
			user := strings.TrimSpace(line[start+len(prefix):])
			if user != "" && !contains(users[user], entry.Name()) {
				users[user] = append(users[user], entry.Name())
			}
		}
		gzr.Close()
		file.Close()
	}

	return users, nil
}

func contains(list []string, value string) bool {
	for _, s := range list {
		if s == value {
			return true
		}
	}
	return false
}
