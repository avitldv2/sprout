package httpd

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Server struct {
	RootDir        string
	LiveReload     *LiveReload
	LivereloadPort int
}

func NewServer(rootDir string) *Server {
	abs, _ := filepath.Abs(rootDir)
	return &Server{
		RootDir: abs,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.LiveReload != nil && r.URL.Path == "/__sprout/events" {
		s.LiveReload.ServeHTTP(w, r)
		return
	}

	path := filepath.Clean(r.URL.Path)
	if strings.Contains(path, "..") {
		http.NotFound(w, r)
		return
	}

	tryPath := func(p string) bool {
		abs, err := filepath.Abs(p)
		if err != nil {
			return false
		}
		rootAbs, _ := filepath.Abs(s.RootDir)
		if !strings.HasPrefix(abs, rootAbs) {
			return false
		}
		if _, err := os.Stat(p); err == nil {
			s.serveFileWithLivereload(w, r, p)
			return true
		}
		return false
	}

	if len(path) > 0 && path[len(path)-1] == '/' {
		if tryPath(filepath.Join(s.RootDir, path, "index.html")) {
			return
		}
	}

	if tryPath(filepath.Join(s.RootDir, path)) {
		return
	}

	if tryPath(filepath.Join(s.RootDir, path, "index.html")) {
		return
	}

	http.NotFound(w, r)
}

func (s *Server) serveFileWithLivereload(w http.ResponseWriter, r *http.Request, filePath string) {
	if !strings.HasSuffix(filePath, ".html") || s.LiveReload == nil {
		http.ServeFile(w, r, filePath)
		return
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var buf strings.Builder
	script := InjectScript(s.LivereloadPort)
	html := string(data)

	if idx := strings.LastIndex(html, "</body>"); idx != -1 {
		buf.WriteString(html[:idx])
		buf.WriteString(script)
		buf.WriteString("\n")
		buf.WriteString(html[idx:])
	} else {
		buf.WriteString(html)
		buf.WriteString(script)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(buf.String()))
}

func Start(port int, rootDir string) error {
	server := NewServer(rootDir)

	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Serving on http://localhost%s\n", addr)

	return http.ListenAndServe(addr, server)
}
