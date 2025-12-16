package app

import (
	"fmt"
	"net/http"
	"path/filepath"
	"sync"
	"time"

	"sprout/internal/cache"
	"sprout/internal/config"
	"sprout/internal/httpd"
	"sprout/internal/logx"
	"sprout/internal/model"
)

type RebuildMode string

const (
	RebuildManual  RebuildMode = "manual"
	RebuildRequest RebuildMode = "request"
	RebuildWatch   RebuildMode = "watch"
)

func Serve(root string, port int, rebuildMode RebuildMode, livereload bool) error {
	cfg, err := config.Load(root)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	resolved, err := config.Resolve(cfg, root)
	if err != nil {
		return fmt.Errorf("failed to resolve config: %w", err)
	}

	if rebuildMode != RebuildManual {
		logx.Infof("Building site...")
		if _, err := Build(root, false); err != nil {
			return fmt.Errorf("initial build failed: %w", err)
		}
	}

	var rebuildFunc func() error
	switch rebuildMode {
	case RebuildManual:
		rebuildFunc = nil
	case RebuildRequest:
		rebuildFunc = createRequestRebuildFunc(root, resolved)
	case RebuildWatch:
		return fmt.Errorf("watch mode not yet implemented")
	default:
		return fmt.Errorf("unknown rebuild mode: %s", rebuildMode)
	}

	var lr *httpd.LiveReload
	if livereload {
		lr = httpd.NewLiveReload()
	}

	baseHandler := httpd.NewServer(resolved.Paths.Public)
	baseHandler.LiveReload = lr
	baseHandler.LivereloadPort = port

	var handler http.Handler = baseHandler
	if rebuildFunc != nil {
		rebuildHandler := &rebuildHandler{
			handler:     baseHandler,
			rebuildFunc: rebuildFunc,
			root:        root,
			resolved:    resolved,
			liveReload:  lr,
		}
		handler = rebuildHandler
	}

	addr := fmt.Sprintf(":%d", port)
	logx.Infof("Serving on http://localhost%d", port)

	return http.ListenAndServe(addr, handler)
}

type rebuildHandler struct {
	handler     http.Handler
	rebuildFunc func() error
	root        string
	resolved    *model.ResolvedConfig
	liveReload  *httpd.LiveReload
	mu          sync.Mutex
	lastRebuild time.Time
	rebuilding  bool
}

func (h *rebuildHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	needsRebuild := false
	now := time.Now()
	if now.Sub(h.lastRebuild) > 500*time.Millisecond && !h.rebuilding {
		needsRebuild = true
		h.rebuilding = true
	}
	h.mu.Unlock()

	if needsRebuild {
		changed, err := checkSourcesChanged(h.root, h.resolved)
		if err != nil {
			logx.Errorf("Failed to check sources: %v", err)
		} else if changed {
			logx.Infof("Sources changed, rebuilding...")
			if err := h.rebuildFunc(); err != nil {
				logx.Errorf("Rebuild failed: %v", err)
			} else if h.liveReload != nil {
				h.liveReload.Reload()
			}
		}

		h.mu.Lock()
		h.lastRebuild = time.Now()
		h.rebuilding = false
		h.mu.Unlock()
	}

	h.handler.ServeHTTP(w, r)
}

func checkSourcesChanged(root string, resolved *model.ResolvedConfig) (bool, error) {
	cachePath := filepath.Join(resolved.Paths.Public, ".sprout-cache.json")
	oldCache, err := cache.LoadCache(cachePath)
	if err != nil {
		return true, nil
	}

	snapshot, err := cache.CreateSnapshot(resolved.Paths.Content, resolved.Paths.Templates, resolved.Paths.Static)
	if err != nil {
		return false, err
	}

	plan := cache.Diff(oldCache, snapshot)
	return len(plan.PagesToRebuild) > 0 || len(plan.AssetsToCopy) > 0 || plan.TemplatesChanged, nil
}

func createRequestRebuildFunc(root string, resolved *model.ResolvedConfig) func() error {
	return func() error {
		_, err := Build(root, false)
		return err
	}
}
