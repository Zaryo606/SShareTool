package main

import (
	"fmt"
	"os"
)

var (
	logFile   *os.File
	stepIndex = 1
	stepTotal = 12
)

func main() {
	fmt.Println("══════════════════════════════════════════")
	fmt.Println("  Dark Minecraft cheat finder")
	fmt.Println("  Basic screensharing tool")
	fmt.Println("  Made by Zary")
	fmt.Println("══════════════════════════════════════════")
	fmt.Println()

	var err error
	logFile, err = os.Create("scanlog.txt")
	if err != nil {
		fmt.Println("Cannot create scanlog.txt:", err)
		return
	}
	defer logFile.Close()

	step("Checking for virtual machine")
	result := checkVM()
	log(result)

	step("Checking Recycle Bin")
	result = checkRecycleBin()
	log(result)

	step("Scanning Prefetch folder")
	result = checkPrefetch()
	log(result)

	step("Checking if explorer was restarted")
	result = checkJavaVsExplorer()
	log(result)

	step("Scanning for exe strings")
	result = checkExplorerStrings()
	log(result)

	step("Checking for jnativehook")
	result = checkTempJNativeHook()
	log(result)

	step("Scanning Recent folder")
	result = checkRecentMacros()
	log(result)

	step("Scanning main disk")
	result = checkSuspiciousNamesOnCDrive()
	log(result)

	step("Checking system services")
	result = checkSystemServices()
	log(result)

	step("Checking Minecraft log ")
	result = checkMinecraftLogs()
	log(result)

	step("Checking Minecraft versions")
	result = checkMinecraftVersions()
	log(result)

	step("Scanning user mods")
	result = checkModJars()
	log(result)

	fmt.Println("\n[~] Scan completed! → Output saved to scanlog.txt")
}

func step(description string) {
	fmt.Printf("Process [%d/%d]: %s\n", stepIndex, stepTotal, description)
	stepIndex++
}

func log(msg string) {
	logFile.WriteString(msg + "\n\n")
}
