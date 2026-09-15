//go:build darwin

// Package macffi seeds a Developer ID–signed libffi into the jupiterrider/ffi
// cache path before that package's init dlopens the library.
//
// jupiterrider/ffi embeds an unsigned libffi.8.dylib and writes it under
// ~/Library/Caches/... Hardened Runtime (AMFI) rejects that load. The signed
// DMG build places a codesigned copy at Contents/Frameworks/libffi.8.dylib;
// this package copies it over the cache file so ffi loads a Team-ID–matched
// dylib. Import this package before any import that pulls in yzma/ffi.
package macffi

import (
	"os"
	"path/filepath"
)

// Must match github.com/jupiterrider/ffi libffiVersion (currently v0.6.0 → 3.5.1).
const libffiCacheVersion = "3.5.1"

func init() {
	exe, err := os.Executable()
	if err != nil {
		return
	}
	seedSignedLibFFI(exe)
}

func seedSignedLibFFI(exe string) {
	src, ok := bundledLibFFIPath(exe)
	if !ok {
		return
	}
	cacheDir, err := os.UserCacheDir()
	if err != nil || cacheDir == "" {
		return
	}
	destDir := filepath.Join(cacheDir, "github.com/jupiterrider/ffi/libffi", libffiCacheVersion)
	dest := filepath.Join(destDir, "libffi.8.dylib")
	data, err := os.ReadFile(src)
	if err != nil || len(data) == 0 {
		return
	}
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return
	}
	tmp, err := os.CreateTemp(destDir, "libffi.8.dylib.*")
	if err != nil {
		return
	}
	tmpName := tmp.Name()
	okWrite := false
	defer func() {
		_ = tmp.Close()
		if !okWrite {
			_ = os.Remove(tmpName)
		}
	}()
	if _, err := tmp.Write(data); err != nil {
		return
	}
	if err := tmp.Chmod(0755); err != nil {
		return
	}
	if err := tmp.Close(); err != nil {
		return
	}
	if err := os.Rename(tmpName, dest); err != nil {
		_ = os.Remove(tmpName)
		return
	}
	okWrite = true
}

func bundledLibFFIPath(exe string) (string, bool) {
	resolved, err := filepath.EvalSymlinks(exe)
	if err == nil {
		exe = resolved
	}
	// Linmo.app/Contents/MacOS/Linmo → ../Frameworks/libffi.8.dylib
	candidate := filepath.Clean(filepath.Join(filepath.Dir(exe), "..", "Frameworks", "libffi.8.dylib"))
	if st, err := os.Stat(candidate); err == nil && !st.IsDir() {
		return candidate, true
	}
	return "", false
}
