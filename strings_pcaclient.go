package main

import (
	"bytes"
	"fmt"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	MEM_COMMIT  = 0x1000
	MEM_PRIVATE = 0x20000
	MEM_IMAGE   = 0x1000000
	MEM_MAPPED  = 0x40000

	ERROR_PARTIAL_COPY syscall.Errno = 299
)

type MEMORY_BASIC_INFORMATION struct {
	BaseAddress       uintptr
	AllocationBase    uintptr
	AllocationProtect uint32
	RegionSize        uintptr
	State             uint32
	Protect           uint32
	Type              uint32
}

func checkExplorerStrings() string {
	const target = "pcaclient"
	const minLength = 4

	pid := findExplorerPID()
	if pid == 0 {
		return "[5] ~ explorer.exe not found"
	}

	handle, err := windows.OpenProcess(windows.PROCESS_QUERY_INFORMATION|windows.PROCESS_VM_READ, false, uint32(pid))
	if err != nil {
		return fmt.Sprintf("[5] ~ Cannot open explorer.exe: %v", err)
	}
	defer windows.CloseHandle(handle)

	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	procVirtualQueryEx := kernel32.NewProc("VirtualQueryEx")

	var addr uintptr = 0
	matches := []string{}

	for {
		var memInfo MEMORY_BASIC_INFORMATION
		memInfoSize := unsafe.Sizeof(memInfo)

		ret, _, _ := procVirtualQueryEx.Call(
			uintptr(handle),
			addr,
			uintptr(unsafe.Pointer(&memInfo)),
			memInfoSize,
		)

		if ret == 0 {
			break
		}

		if memInfo.State == MEM_COMMIT && (memInfo.Type == MEM_PRIVATE || memInfo.Type == MEM_IMAGE || memInfo.Type == MEM_MAPPED) {
			size := int(memInfo.RegionSize)
			if size == 0 || size > 100*1024*1024 {
				addr += memInfo.RegionSize
				continue
			}

			buffer := make([]byte, size)
			var bytesRead uintptr

			errRead := windows.ReadProcessMemory(handle, memInfo.BaseAddress, &buffer[0], uintptr(size), &bytesRead)
			if errRead != nil && errRead != ERROR_PARTIAL_COPY {
				addr += memInfo.RegionSize
				continue
			}

			found := fastSearchStrings(buffer[:bytesRead], target, minLength)
			matches = append(matches, found...)
		}

		addr += memInfo.RegionSize
	}

	if len(matches) == 0 {
		return "[5] ~ No matches for 'pcaclient' found in explorer.exe"
	}

	var builder strings.Builder
	builder.WriteString("[5] ~ Detected 'pcaclient' strings in explorer.exe:\n")

	for _, match := range matches {
		lines := strings.Split(match, "\n")
		for i, line := range lines {
			if i == 0 {
				builder.WriteString("      - " + line + "\n")
			} else {
				builder.WriteString("        " + line + "\n")
			}
		}
	}
	return builder.String()
}

func findExplorerPID() int {
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return 0
	}
	defer windows.CloseHandle(snapshot)

	var entry windows.ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))

	if err := windows.Process32First(snapshot, &entry); err != nil {
		return 0
	}

	for {
		name := syscall.UTF16ToString(entry.ExeFile[:])
		if strings.EqualFold(name, "explorer.exe") {
			return int(entry.ProcessID)
		}
		if err := windows.Process32Next(snapshot, &entry); err != nil {
			break
		}
	}
	return 0
}

func fastSearchStrings(buffer []byte, needle string, minLength int) []string {
	var results []string
	needleLower := []byte(strings.ToLower(needle))
	bufferLower := bytes.ToLower(buffer)

	if !bytes.Contains(bufferLower, needleLower) {
		return results
	}

	for _, part := range bytes.Split(buffer, []byte{0}) {
		if len(part) >= minLength && bytes.Contains(bytes.ToLower(part), needleLower) {
			results = append(results, string(part))
		}
	}
	return results
}
