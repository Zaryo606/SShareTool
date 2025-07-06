package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/process"
)

func checkMinecraftLogs() string {
	appdata := os.Getenv("APPDATA")
	logPath := filepath.Join(appdata, ".minecraft", "logs", "latest.log")
	logsDir := filepath.Dir(logPath)

	// Najdi běžící javaw.exe
	startTime, err := getJavawStartTime()
	if err != nil {
		return "[10] ~ Minecraft isnt running, or its using other service"
	}

	// Zkontroluj, zda byl latest.log změněn po startu MC
	info, err := os.Stat(logPath)
	if err != nil {
		return "[10] ~ Couldt fetch latest log"
	}
	if info.ModTime().Before(startTime) {
		return "[10] ~ Minecraft logging system isnt working"
	}

	// Zjisti aktivního uživatele z latest.log
	activeUser, err := getLatestUserFromLog(logPath)
	if err != nil {
		return "[10] ~ Couldt find active user"
	}

	// Získání unikátních jmen z archivovaných logů
	userList, err := getUniqueUsersFromGzLogs(logsDir)
	if err != nil {
		return "[10] ~ Error handling older logs: " + err.Error()
	}

	var builder strings.Builder
	builder.WriteString("[10] ~ Minecraft log check:\n")
	builder.WriteString("      - Active user: " + activeUser + "\n")
	builder.WriteString("      - Other usernames: \n")
	for _, user := range userList {
		builder.WriteString("         • " + user + "\n")
	}

	return builder.String()
}

func getJavawStartTime() (time.Time, error) {
	procs, err := process.Processes()
	if err != nil {
		return time.Time{}, err
	}
	for _, p := range procs {
		name, err := p.Name()
		if err == nil && strings.EqualFold(name, "javaw.exe") {
			createTime, err := p.CreateTime()
			if err != nil {
				return time.Time{}, err
			}
			return time.UnixMilli(createTime), nil
		}
	}
	return time.Time{}, fmt.Errorf("javaw.exe nenalezen")
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
		prefix := "Successfully refreshed token for "
		if strings.Contains(line, prefix) {
			start := strings.Index(line, prefix) + len(prefix)
			if start < len(line) {
				lastUser = strings.TrimSpace(line[start:])
			}
		}
	}
	if lastUser == "" {
		return "", fmt.Errorf("uživatel nenalezen")
	}
	return lastUser, nil
}

func getUniqueUsersFromGzLogs(logsDir string) ([]string, error) {
	userSet := make(map[string]struct{})

	err := filepath.Walk(logsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".gz") {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		gzReader, err := gzip.NewReader(f)
		if err != nil {
			return err
		}
		defer gzReader.Close()

		scanner := bufio.NewScanner(gzReader)
		for scanner.Scan() {
			line := scanner.Text()
			prefix := "Successfully refreshed token for "
			if strings.Contains(line, prefix) {
				start := strings.Index(line, prefix) + len(prefix)
				if start < len(line) {
					name := strings.TrimSpace(line[start:])
					userSet[name] = struct{}{}
				}
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	users := make([]string, 0, len(userSet))
	for name := range userSet {
		users = append(users, name)
	}
	sort.Strings(users)
	return users, nil
}
