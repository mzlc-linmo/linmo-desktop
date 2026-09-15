//go:build darwin

package macffi

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSeedSignedLibFFI(t *testing.T) {
	tmp := t.TempDir()
	frameworks := filepath.Join(tmp, "Contents", "Frameworks")
	macos := filepath.Join(tmp, "Contents", "MacOS")
	if err := os.MkdirAll(frameworks, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(macos, 0755); err != nil {
		t.Fatal(err)
	}
	src := filepath.Join(frameworks, "libffi.8.dylib")
	payload := []byte("fake-signed-libffi")
	if err := os.WriteFile(src, payload, 0755); err != nil {
		t.Fatal(err)
	}
	exe := filepath.Join(macos, "Linmo")
	if err := os.WriteFile(exe, []byte("#!/bin/sh\n"), 0755); err != nil {
		t.Fatal(err)
	}

	home := t.TempDir()
	t.Setenv("HOME", home)

	seedSignedLibFFI(exe)

	want := filepath.Join(home, "Library", "Caches", "github.com/jupiterrider/ffi/libffi", libffiCacheVersion, "libffi.8.dylib")
	got, err := os.ReadFile(want)
	if err != nil {
		t.Fatalf("read seeded libffi: %v", err)
	}
	if string(got) != string(payload) {
		t.Fatalf("seeded libffi = %q, want %q", got, payload)
	}
}

func TestSeedSignedLibFFINoopWithoutBundle(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	exe := filepath.Join(t.TempDir(), "not-an-app-binary")
	if err := os.WriteFile(exe, []byte("x"), 0755); err != nil {
		t.Fatal(err)
	}
	seedSignedLibFFI(exe)
	want := filepath.Join(home, "Library", "Caches", "github.com/jupiterrider/ffi/libffi", libffiCacheVersion, "libffi.8.dylib")
	if _, err := os.Stat(want); !os.IsNotExist(err) {
		t.Fatalf("expected no cache file, stat err=%v", err)
	}
}
