package connection

import (
    "path/filepath"
	"strings"
)

func isExecutable(filename string) bool {
	ext := filepath.Ext(filename)
	return ext == ".exe"
}

func ConvertLinuxToWindowsPath(linuxPath string) string {
	linuxPath = strings.ReplaceAll(linuxPath, "/", "\\")

	if len(linuxPath) >= 4 && strings.HasPrefix(linuxPath, "/mnt/") {
		driveLetter := strings.ToUpper(linuxPath[5:6])
		return driveLetter + ":" + linuxPath[6:]
	}

	return linuxPath
}

func ConvertWindowsToLinuxPath(windowsPath string) string {
	windowsPath = strings.ReplaceAll(windowsPath, "\\", "/")

	if len(windowsPath) >= 3 && windowsPath[1] == ':' && windowsPath[2] == '/' {
		driveLetter := windowsPath[0]
		return "/mnt/" + strings.ToLower(string(driveLetter)) + windowsPath[2:]
	}

	return windowsPath
}