// Package logging settings: log dir resolution and default rolling path.
// Mirrors src/logging/logger.ts and src/infra/tmp-linmo-dir.ts.
package logging

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	// LogPrefix is the rolling log file prefix (linmo-YYYY-MM-DD.log).
	LogPrefix = "linmo"
	// LogSuffix is the rolling log file suffix.
	LogSuffix = ".log"
	// PosixOpenClawTmpDir is the preferred log directory on POSIX (same as Node).
	PosixOpenClawTmpDir = "/tmp/linmo"
)

// ResolvePreferredLogDir returns the preferred directory for log files.
// Same logic as Node: /tmp/linmo if it exists and is writable, else os.TempDir()/linmo.
func ResolvePreferredLogDir() string {
	info, err := os.Stat(PosixOpenClawTmpDir)
	if err != nil {
		if os.IsNotExist(err) {
			return filepath.Join(os.TempDir(), "linmo")
		}
		return filepath.Join(os.TempDir(), "linmo")
	}
	if !info.IsDir() {
		return filepath.Join(os.TempDir(), "linmo")
	}
	// Check writable
	f, err := os.Create(filepath.Join(PosixOpenClawTmpDir, ".write-check"))
	if err != nil {
		return filepath.Join(os.TempDir(), "linmo")
	}
	f.Close()
	_ = os.Remove(f.Name())
	return PosixOpenClawTmpDir
}

// DefaultRollingLogPath returns the default log file path for today: linmo-YYYY-MM-DD.log.
// Mirrors Node defaultRollingPathForToday().
func DefaultRollingLogPath(logDir string) string {
	if logDir == "" {
		logDir = ResolvePreferredLogDir()
	}
	today := time.Now().UTC().Format("2006-01-02")
	return filepath.Join(logDir, LogPrefix+"-"+today+LogSuffix)
}

// ResolvedSettings holds the resolved log file path (for logs.tail and file logger).
type ResolvedSettings struct {
	File string
}

// GetResolvedLoggerSettings returns the log file path.
// Uses config.logging.file if available via env/config, else default rolling path in preferred dir.
func GetResolvedLoggerSettings(env func(string) string, logDirOverride string) ResolvedSettings {
	if env == nil {
		env = func(k string) string { return os.Getenv(k) }
	}
	// Optional: LIMNO_LOG_DIR to force log dir (e.g. state dir)
	logDir := strings.TrimSpace(logDirOverride)
	if logDir == "" {
		logDir = strings.TrimSpace(env("LIMNO_LOG_DIR"))
	}
	if logDir != "" {
		return ResolvedSettings{File: DefaultRollingLogPath(logDir)}
	}
	return ResolvedSettings{File: DefaultRollingLogPath(ResolvePreferredLogDir())}
}

// RuntimeVersion returns Go runtime version (e.g. "go1.21.0") for _meta.runtimeVersion.
func RuntimeVersion() string {
	return strings.TrimPrefix(runtime.Version(), "go")
}
