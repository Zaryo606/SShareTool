package main

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func wmiStartTime(processName string) (time.Time, error) {
	out, err := exec.Command("wmic", "process", "where", fmt.Sprintf("name='%s'", processName), "get", "CreationDate", "/value").Output()
	if err != nil {
		return time.Time{}, err
	}
	for _, line := range strings.Split(string(out), "\n") {
		if strings.HasPrefix(line, "CreationDate=") {
			val := strings.TrimPrefix(line, "CreationDate=")
			if len(val) >= 14 {
				return time.Parse("20060102150405", val[:14])
			}
		}
	}
	return time.Time{}, errors.New("creation date not found")
}

func checkJavaVsExplorer() string {
	javaTime, javaErr := wmiStartTime("javaw.exe")
	explTime, explErr := wmiStartTime("explorer.exe")

	if javaErr != nil {
		return "[4] ~ javaw.exe not running"
	}
	if explErr != nil {
		return "[4] ~ explorer.exe not running"
	}

	now := time.Now()
	javaUp := now.Sub(javaTime).Round(time.Second)
	explUp := now.Sub(explTime).Round(time.Second)

	result := fmt.Sprintf(
		"[4] ~ Process runtime:\n"+
			"      - javaw.exe:     since %s (%s)\n"+
			"      - explorer.exe:  since %s (%s)\n",
		javaTime.Format("02.01.2006 15:04:05"), javaUp,
		explTime.Format("02.01.2006 15:04:05"), explUp,
	)

	if javaTime.Before(explTime) {
		result += "      - Result: javaw.exe started BEFORE explorer.exe"
	} else {
		result += "      - Result: javaw.exe started AFTER explorer.exe"
	}

	return result
}
