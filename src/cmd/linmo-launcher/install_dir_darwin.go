//go:build darwin

package main

import (
	"os"
	"path/filepath"
)

// resolveInstallStateDir returns the state directory when running from .app bundle.
// Uses Contents/Resources/data relative to the executable.
func resolveInstallStateDir() string {
	self, err := os.Executable()
	if err != nil {
		return ""
	}
	// self = .../Linmo.app/Contents/MacOS/Linmo
	macosDir := filepath.Dir(self)
	contentsDir := filepath.Join(macosDir, "..")
	return filepath.Join(contentsDir, "Resources", "data")
}
