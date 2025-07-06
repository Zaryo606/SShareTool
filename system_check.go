package main

import (
	"fmt"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

func checkSystemServices() string {
	// Kritické služby
	targets := map[string]string{
		"SysMain":   "Superfetch",
		"PcaSvc":    "Program Compatibility Assistant",
		"WSearch":   "Windows Search",
		"mpssvc":    "Windows Defender Firewall",
		"WdNisSvc":  "Windows Defender Network Inspection Service",
		"WinDefend": "Windows Defender Antivirus",
		"wscsvc":    "Windows Security Center",
	}

	// Výsledky
	results := map[string]string{} // mapa: název služby -> status ("running", "stopped", "missing")
	for name := range targets {
		results[name] = "missing"
	}

	// Připojit se k Service Control Manager
	handle, err := windows.OpenSCManager(nil, nil, windows.SC_MANAGER_ENUMERATE_SERVICE)
	if err != nil {
		return fmt.Sprintf("[9] ~ Failed to open Service Manager: %v", err)
	}
	defer windows.CloseServiceHandle(handle)

	var needed, returned, resume uint32
	var buf [1 << 16]byte

	err = windows.EnumServicesStatusEx(
		handle,
		windows.SC_ENUM_PROCESS_INFO,
		windows.SERVICE_WIN32,
		windows.SERVICE_STATE_ALL,
		(*byte)(unsafe.Pointer(&buf[0])),
		uint32(len(buf)),
		&needed,
		&returned,
		&resume,
		nil,
	)
	if err != nil {
		return fmt.Sprintf("[9] ~ Failed to enumerate services: %v", err)
	}

	entries := (*[1 << 10]windows.ENUM_SERVICE_STATUS_PROCESS)(unsafe.Pointer(&buf[0]))[:returned:returned]
	for _, svc := range entries {
		name := windows.UTF16PtrToString(svc.ServiceName)
		if _, ok := results[name]; ok {
			if svc.ServiceStatusProcess.CurrentState == windows.SERVICE_RUNNING {
				results[name] = "running"
			} else {
				results[name] = "stopped"
			}
		}
	}

	// Výstup
	var sus []string
	for name, state := range results {
		switch state {
		case "missing":
			sus = append(sus, fmt.Sprintf("      - Service '%s' (%s) is MISSING", name, targets[name]))
		case "stopped":
			sus = append(sus, fmt.Sprintf("      - Service '%s' (%s) is NOT running", name, targets[name]))
		}
	}

	if len(sus) == 0 {
		return "[9] ~ All critical system services are present and running"
	}

	var builder strings.Builder
	builder.WriteString("[9] ~ Suspicious system service state:\n")
	for _, line := range sus {
		builder.WriteString(line + "\n")
	}
	return builder.String()
}
