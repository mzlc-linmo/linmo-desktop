package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/openocta/openocta/pkg/appinstance"
)

const (
	defaultPort = 18900
)

func main() {
	appinstance.KillOtherLinmoProcesses()

	port := defaultPort
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	baseURL := fmt.Sprintf("http://%s", addr)

	selfDir, _ := os.Executable()
	binDir := filepath.Dir(selfDir)

	linmoPath := filepath.Join(binDir, binaryName("linmo"))
	if _, err := os.Stat(linmoPath); err != nil {
		// Fallback: rely on PATH
		linmoPath = binaryName("linmo")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start gateway in desktop mode.
	cmd := exec.CommandContext(ctx, linmoPath, "gateway", "run", "--port", fmt.Sprintf("%d", port))
	cmd.Env = append(os.Environ(), "LIMNO_RUN_MODE=desktop")
	if stateDir := resolveInstallStateDir(); stateDir != "" {
		_ = os.MkdirAll(stateDir, 0755)
		cmd.Env = append(cmd.Env, "LIMNO_STATE_DIR="+stateDir)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = binDir

	_ = cmd.Start()

	// Wait for health.
	_ = waitForHealthy(baseURL+"/health", 12*time.Second)

	// Open browser.
	_ = openBrowser(baseURL)

	// Keep launcher alive while gateway is alive (best-effort).
	_ = cmd.Wait()
}

func binaryName(base string) string {
	if runtime.GOOS == "windows" {
		return base + ".exe"
	}
	return base
}

func waitForHealthy(url string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	client := &http.Client{Timeout: 800 * time.Millisecond}
	for time.Now().Before(deadline) {
		req, _ := http.NewRequest("GET", url, nil)
		res, err := client.Do(req)
		if err == nil {
			_ = res.Body.Close()
			if res.StatusCode >= 200 && res.StatusCode < 500 {
				return nil
			}
		}
		time.Sleep(250 * time.Millisecond)
	}
	return errors.New("gateway health check timeout")
}

func openBrowser(url string) error {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", url).Start()
	case "windows":
		// Use cmd start to open the default browser.
		return exec.Command("cmd", "/c", "start", "", url).Start()
	default:
		// linux and others
		if _, err := exec.LookPath("xdg-open"); err == nil {
			return exec.Command("xdg-open", url).Start()
		}
	}
	return nil
}
