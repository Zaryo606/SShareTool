package main

import (
	"fmt"
	"os/exec"
	"strings"
)

var vmIndicators = []string{
	"vmware", "virtualbox", "kvm", "qemu", "virtual machine",
	"parallels", "hyper-v", "xen", "vbox", "bochs",
}

func getSystemModel() (manufacturer string, model string, err error) {
	out, err := exec.Command("wmic", "computersystem", "get", "manufacturer,", "model").Output()
	if err != nil {
		return "", "", err
	}

	lines := strings.Split(string(out), "\n")
	if len(lines) < 2 {
		return "", "", fmt.Errorf("output too short")
	}

	data := strings.Fields(lines[1])
	if len(data) < 2 {
		return "", "", fmt.Errorf("could not parse model info")
	}

	return data[0], strings.Join(data[1:], " "), nil
}

func checkVM() string {
	manufacturer, model, err := getSystemModel()
	if err != nil {
		return fmt.Sprintf("[1] ~ Failed to detect system model: %v", err)
	}

	lowerModel := strings.ToLower(model)
	lowerManuf := strings.ToLower(manufacturer)

	for _, indicator := range vmIndicators {
		if strings.Contains(lowerModel, indicator) || strings.Contains(lowerManuf, indicator) {
			return fmt.Sprintf("[1] ~ Likely running in VM!\n     -> Manufacturer: %s\n     -> Model: %s", manufacturer, model)
		}
	}

	return fmt.Sprintf("[1] ~ Physical machine detected\n     -> Manufacturer: %s\n     -> Model: %s", manufacturer, model)
}
