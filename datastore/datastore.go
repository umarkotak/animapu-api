package datastore

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
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
	Db      *sqlx.DB // required
	GoCache *cache.Cache
	Browser *rod.Browser
}

var dataStore DataStore

func Initialize() error {
	db, err := sqlx.Connect("postgres", config.Get().DbUrl)
	if err != nil {
		logrus.Error(err)
		return err
	}

	goCache := cache.New(5*time.Minute, 10*time.Minute)
	browser, err := launchChrome()
	if err != nil {
		return err
	}

	dataStore = DataStore{
		Db:      db,
		GoCache: goCache,
		Browser: browser,
	}

	return nil
}

func launchChrome() (*rod.Browser, error) {
	var version struct {
		WebSocketDebuggerURL string `json:"webSocketDebuggerUrl"`
	}
	client := http.Client{Timeout: time.Second}
	closedExistingChrome := false

	resp, err := client.Get("http://127.0.0.1:9222/json/version")
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
				resp, err := client.Get("http://127.0.0.1:9222/json/version")
				if err != nil {
					break
				}
				resp.Body.Close()
				time.Sleep(250 * time.Millisecond)
			}
			if closedExistingChrome {
				resp, err := client.Get("http://127.0.0.1:9222/json/version")
				if err == nil {
					resp.Body.Close()
					return nil, fmt.Errorf("existing Chrome did not close")
				}
			}
			args := []string{"-na", "Google Chrome", "--args", "--remote-debugging-port=9222", "--user-data-dir=/tmp/chrome-rod", "--mute-audio"}
			if config.Get().RodHeadless {
				args = append(args, "--headless=new")
			}
			if err := exec.Command("open", args...).Run(); err != nil {
				return nil, fmt.Errorf("launch Chrome: %w", err)
			}
		}

		resp, err := client.Get("http://127.0.0.1:9222/json/version")
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

	return nil, fmt.Errorf("Chrome debugger did not start")
}

func Close() {
	if dataStore.Browser != nil {
		if err := dataStore.Browser.Close(); err != nil {
			logrus.WithError(err).Warn("close Chrome")
		}
	}
}

// NewBrowser creates an isolated browser context using the Chrome instance started by Initialize.
func NewBrowser() (*rod.Browser, error) {
	var version struct {
		WebSocketDebuggerURL string `json:"webSocketDebuggerUrl"`
	}

	resp, err := (&http.Client{Timeout: time.Second}).Get("http://127.0.0.1:9222/json/version")
	if err != nil {
		return nil, fmt.Errorf("get Chrome debugger URL: %w", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&version); err != nil {
		return nil, fmt.Errorf("decode Chrome debugger URL: %w", err)
	}
	if version.WebSocketDebuggerURL == "" {
		return nil, fmt.Errorf("Chrome debugger URL is empty")
	}

	browser := rod.New().ControlURL(version.WebSocketDebuggerURL)
	if err := browser.Connect(); err != nil {
		return nil, fmt.Errorf("connect to Chrome: %w", err)
	}

	isolatedBrowser, err := browser.Incognito()
	if err != nil {
		return nil, fmt.Errorf("create isolated browser context: %w", err)
	}
	return isolatedBrowser, nil
}

func Get() DataStore {
	return dataStore
}
