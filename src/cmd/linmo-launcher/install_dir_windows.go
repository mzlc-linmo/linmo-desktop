//go:build windows

package main

import (
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

const regPath = `SOFTWARE\Linmo`
const regValue = "InstallDir"

// resolveInstallStateDir returns the state directory when installed via NSIS.
// Reads InstallDir from HKLM and returns InstallDir\data.
func resolveInstallStateDir() string {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, regPath, registry.QUERY_VALUE)
	if err != nil {
		return ""
	}
	defer k.Close()
	installDir, _, err := k.GetStringValue(regValue)
	if err != nil || installDir == "" {
		return ""
	}
	return filepath.Join(installDir, "data")
}
