package datastore

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
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
var browsersMu sync.Mutex

func Initialize() error {
	db, err := sqlx.Connect("postgres", config.Get().DbUrl)
	if err != nil {
		logrus.Error(err)
		return err
	}

	goCache := cache.New(5*time.Minute, 10*time.Minute)
	poolSize := config.Get().RodBrowserPoolSize
	pool := make(chan *rod.Browser, poolSize)

	dataStore = DataStore{
		Db:          db,
		GoCache:     goCache,
		BrowserPool: pool,
	}
	for slot := 0; slot < poolSize; slot++ {
		go retryChrome(pool, 9222+slot)
	}

	return nil
}

func retryChrome(pool chan<- *rod.Browser, port int) {
	for {
		browser, err := launchChrome(port)
		if err == nil {
			browsersMu.Lock()
			dataStore.Browsers = append(dataStore.Browsers, browser)
			browsersMu.Unlock()
			pool <- browser
			return
		}

		logrus.WithError(err).Warnf("launch Chrome on port %d; retrying", port)
		time.Sleep(5 * time.Second)
	}
}

func launchChrome(port int) (*rod.Browser, error) {
	bin, exists := launcher.LookPath()
	if !exists {
		return nil, fmt.Errorf("Chrome executable not found")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Disabling leakless lets Rod reuse an existing debugger on this port.
	u, err := launcher.New().
		Bin(bin).
		Context(ctx).
		Leakless(false).
		RemoteDebuggingPort(port).
		UserDataDir(filepath.Join(os.TempDir(), fmt.Sprintf("chrome-rod-%d", port))).
		HeadlessNew(config.Get().RodHeadless).
		Set("mute-audio").
		Launch()
	if err != nil {
		return nil, fmt.Errorf("launch Chrome on port %d: %w", port, err)
	}

	browser := rod.New().ControlURL(u)
	if err := browser.Connect(); err != nil {
		return nil, fmt.Errorf("connect to Chrome on port %d: %w", port, err)
	}
	return browser, nil
}

func Close() {
	browsersMu.Lock()
	browsers := dataStore.Browsers
	browsersMu.Unlock()
	for _, browser := range browsers {
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
