package datastore

import (
	"net"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/umarkotak/animapu-api/config"
)

func TestLaunchChromeStartsAndReusesBrowser(t *testing.T) {
	if os.Getenv("RUN_CHROME_TESTS") != "1" {
		t.Skip("set RUN_CHROME_TESTS=1 to test with installed Chrome")
	}
	t.Setenv("ROD_HEADLESS", "true")
	t.Setenv("TMPDIR", t.TempDir())
	if err := config.Initialize(); err != nil {
		t.Fatal(err)
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	if err := listener.Close(); err != nil {
		t.Fatal(err)
	}

	browser, err := launchChrome(port)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = browser.Close() })
	processes, err := (proto.SystemInfoGetProcessInfo{}).Call(browser)
	if err != nil {
		t.Fatal(err)
	}
	for _, info := range processes.ProcessInfo {
		if info.Type != "browser" {
			continue
		}
		process, err := os.FindProcess(info.ID)
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			_ = browser.Close()
			deadline := time.Now().Add(5 * time.Second)
			for process.Signal(syscall.Signal(0)) == nil {
				if time.Now().After(deadline) {
					t.Error("Chrome did not exit before profile cleanup")
					return
				}
				time.Sleep(50 * time.Millisecond)
			}
		})
		break
	}
	url, err := launcher.ResolveURL(listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	reused, err := launchChrome(port)
	if err != nil {
		t.Fatal(err)
	}
	reusedURL, err := launcher.ResolveURL(listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	if reusedURL != url {
		t.Fatalf("Chrome restarted: debugger URL changed from %q to %q", url, reusedURL)
	}
	// Both connections must remain usable after the launch timeout is canceled.
	if _, err := (proto.BrowserGetVersion{}).Call(browser); err != nil {
		t.Fatal(err)
	}
	if _, err := (proto.BrowserGetVersion{}).Call(reused); err != nil {
		t.Fatal(err)
	}
}
