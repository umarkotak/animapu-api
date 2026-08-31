package datastore

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"sync"
	"time"

	"github.com/go-rod/rod"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/config"
)

type DataStore struct {
	Db          *sqlx.DB // required
	GoCache     *cache.Cache
	BrowserPool chan *rod.Browser
	Browsers    []*rod.Browser
}

type BrowserLease struct {
	Browser *rod.Browser
	root    *rod.Browser
	pool    chan *rod.Browser
	once    sync.Once
}

var dataStore DataStore

func Initialize() error {
	db, err := sqlx.Connect("postgres", config.Get().DbUrl)
	if err != nil {
		logrus.Error(err)
		return err
	}

	goCache := cache.New(5*time.Minute, 10*time.Minute)
	poolSize := config.Get().RodBrowserPoolSize
	browsers := make([]*rod.Browser, 0, poolSize)
	pool := make(chan *rod.Browser, poolSize)
	for slot := 0; slot < poolSize; slot++ {
		browser, err := launchChrome(9222 + slot)
		if err != nil {
			for _, startedBrowser := range browsers {
				_ = startedBrowser.Close()
			}
			return err
		}
		browsers = append(browsers, browser)
		pool <- browser
	}

	dataStore = DataStore{
		Db:          db,
		GoCache:     goCache,
		BrowserPool: pool,
		Browsers:    browsers,
	}

	return nil
}

func launchChrome(port int) (*rod.Browser, error) {
	var version struct {
		WebSocketDebuggerURL string `json:"webSocketDebuggerUrl"`
	}
	client := http.Client{Timeout: time.Second}
	closedExistingChrome := false
	debuggerURL := fmt.Sprintf("http://127.0.0.1:%d/json/version", port)

	resp, err := client.Get(debuggerURL)
	if err == nil {
		err = json.NewDecoder(resp.Body).Decode(&version)
		resp.Body.Close()
		if err == nil && version.WebSocketDebuggerURL != "" {
			browser := rod.New().ControlURL(version.WebSocketDebuggerURL)
			if err := browser.Connect(); err != nil {
				return nil, fmt.Errorf("connect to existing Chrome: %w", err)
			}
			if err := browser.Close(); err != nil {
				return nil, fmt.Errorf("close existing Chrome: %w", err)
			}
			closedExistingChrome = true
		}
	}

	for attempt := 0; attempt <= 20; attempt++ {
		if attempt == 0 {
			for waitAttempt := 0; closedExistingChrome && waitAttempt <= 20; waitAttempt++ {
				resp, err := client.Get(debuggerURL)
				if err != nil {
					break
				}
				resp.Body.Close()
				time.Sleep(250 * time.Millisecond)
			}
			if closedExistingChrome {
				resp, err := client.Get(debuggerURL)
				if err == nil {
					resp.Body.Close()
					return nil, fmt.Errorf("existing Chrome did not close")
				}
			}
			args := []string{"-na", "Google Chrome", "--args", fmt.Sprintf("--remote-debugging-port=%d", port), fmt.Sprintf("--user-data-dir=/tmp/chrome-rod-%d", port), "--mute-audio"}
			if config.Get().RodHeadless {
				args = append(args, "--headless=new")
			}
			if err := exec.Command("open", args...).Run(); err != nil {
				return nil, fmt.Errorf("launch Chrome: %w", err)
			}
		}

		resp, err := client.Get(debuggerURL)
		if err == nil {
			err = json.NewDecoder(resp.Body).Decode(&version)
			resp.Body.Close()
			if err == nil && version.WebSocketDebuggerURL != "" {
				browser := rod.New().ControlURL(version.WebSocketDebuggerURL)
				if err := browser.Connect(); err != nil {
					return nil, fmt.Errorf("connect to Chrome: %w", err)
				}
				return browser, nil
			}
		}
		time.Sleep(250 * time.Millisecond)
	}

	return nil, fmt.Errorf("Chrome debugger on port %d did not start", port)
}

func Close() {
	for _, browser := range dataStore.Browsers {
		if err := browser.Close(); err != nil {
			logrus.WithError(err).Warn("close Chrome")
		}
	}
}

// NewBrowser borrows an isolated browser context from the Chrome pool.
func NewBrowser(ctx context.Context) (*BrowserLease, error) {
	if dataStore.BrowserPool == nil {
		return nil, fmt.Errorf("browser pool is not initialized")
	}

	select {
	case root := <-dataStore.BrowserPool:
		browser, err := root.Context(ctx).Incognito()
		if err != nil {
			dataStore.BrowserPool <- root
			return nil, fmt.Errorf("create isolated browser context: %w", err)
		}
		return &BrowserLease{Browser: browser, root: root, pool: dataStore.BrowserPool}, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (lease *BrowserLease) Release() {
	lease.once.Do(func() {
		if err := lease.Browser.Context(context.Background()).Close(); err != nil {
			logrus.WithError(err).Warn("close isolated browser context")
		}
		lease.pool <- lease.root
	})
}

func Get() DataStore {
	return dataStore
}
