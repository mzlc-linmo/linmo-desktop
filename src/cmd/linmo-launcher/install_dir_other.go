//go:build !windows && !darwin

package main

func resolveInstallStateDir() string {
	return ""
}
