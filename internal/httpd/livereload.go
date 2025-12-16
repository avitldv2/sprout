package httpd

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type LiveReload struct {
	clients map[chan string]bool
	mu      sync.Mutex
}

func NewLiveReload() *LiveReload {
	return &LiveReload{
		clients: make(map[chan string]bool),
	}
}

func (lr *LiveReload) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	ch := make(chan string, 1)

	lr.mu.Lock()
	lr.clients[ch] = true
	lr.mu.Unlock()

	defer func() {
		lr.mu.Lock()
		delete(lr.clients, ch)
		close(ch)
		lr.mu.Unlock()
	}()

	fmt.Fprintf(w, "data: {\"type\":\"connected\"}\n\n")
	w.(http.Flusher).Flush()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case event := <-ch:
			fmt.Fprintf(w, "data: %s\n\n", event)
			w.(http.Flusher).Flush()
		case <-ticker.C:
			io.WriteString(w, ": keepalive\n\n")
			w.(http.Flusher).Flush()
		case <-r.Context().Done():
			return
		}
	}
}

func (lr *LiveReload) Reload() {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	event := `{"type":"reload"}`
	for ch := range lr.clients {
		select {
		case ch <- event:
		default:
		}
	}
}

func InjectScript(port int) string {
	return fmt.Sprintf(`<script>
(function() {
  var es = new EventSource('http://localhost:%d/__sprout/events');
  es.onmessage = function(e) {
    var data = JSON.parse(e.data);
    if (data.type === 'reload') {
      window.location.reload();
    }
  };
})();
</script>`, port)
}
