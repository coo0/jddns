package common

import (
	"os"
	"path/filepath"
	"runtime"
)

// Determine whether the current system is a Windows system?
func IsWindows() bool {
	if runtime.GOOS == "windows" {
		return true
	}
	return false
}
func GetInstallPath() string {
	var path string
	if IsWindows() {
		path = `C:\Program Files\jddns`
	} else {
		path = "/etc/jddns"
	}
	return path
}
func GetAppPath() string {
	if path, err := filepath.Abs(filepath.Dir(os.Args[0])); err == nil {
		return path
	}
	return os.Args[0]
}

func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func GetLogPath() string {
	var path string
	if IsWindows() {
		path = filepath.Join(GetAppPath(), "jddns.log")
	} else {
		path = "/var/log/jddns.log"
	}
	return path
}
