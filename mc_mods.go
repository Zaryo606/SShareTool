package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func checkModJars() string {
	appdata := os.Getenv("APPDATA")
	modsPath := filepath.Join(appdata, ".minecraft", "mods")

	files, err := os.ReadDir(modsPath)
	if err != nil {
		return "[12] ~ Could load mods folder:  " + err.Error()
	}

	var results []string
	suspicious := []string{
		"aimbot", "killaura", "raven", "doomsday",
		"clicker", "autoclicker", "leftclicker", "antikb", "aristois",
		"aimassist", "autojump", "haru", "wurst", "impact", "sigma",
	}

	for _, file := range files {
		if !strings.HasSuffix(strings.ToLower(file.Name()), ".jar") {
			continue
		}
		fullPath := filepath.Join(modsPath, file.Name())
		matches, err := scanJar(fullPath, suspicious)
		if err != nil {
			results = append(results, fmt.Sprintf("      - %s → cant load (%v)", file.Name(), err))
			continue
		}
		if len(matches) > 0 {
			results = append(results, fmt.Sprintf("      - %s:", file.Name()))
			for _, m := range matches {
				results = append(results, "           → "+m)
			}
		}
	}

	if len(results) == 0 {
		return "[12] ~ Didnt find any suspicous mods in mods folder"
	}

	var builder strings.Builder
	builder.WriteString("[12] ~ Found suspicous mods:\n")
	for _, line := range results {
		builder.WriteString(line + "\n")
	}
	return builder.String()
}

func scanJar(path string, keywords []string) ([]string, error) {
	var found []string
	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	lowerPath := strings.ToLower(path)
	for _, keyword := range keywords {
		if strings.Contains(lowerPath, keyword) {
			found = append(found, fmt.Sprintf("file name includes '%s'", keyword))
		}
	}

	for _, f := range r.File {
		lowerName := strings.ToLower(f.Name)
		for _, keyword := range keywords {
			if strings.Contains(lowerName, keyword) {
				found = append(found, fmt.Sprintf("class/dir name includes '%s' (%s)", keyword, f.Name))
				break
			}
		}

		if strings.HasSuffix(f.Name, ".class") {
			rc, err := f.Open()
			if err != nil {
				continue
			}
			data, _ := io.ReadAll(rc)
			rc.Close()

			content := strings.ToLower(string(data))
			for _, keyword := range keywords {
				if strings.Contains(content, keyword) {
					found = append(found, fmt.Sprintf(".class include '%s' (%s)", keyword, f.Name))
					break
				}
			}
		}
	}

	return found, nil
}
