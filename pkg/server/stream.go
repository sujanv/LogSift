package server

import (
	"fmt"
	"net/http"
)

// HandleStream provides Server-Sent Events (SSE) live streaming of log events.
func (s *Server) HandleStream(lineChan <-chan []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		for {
			select {
			case <-r.Context().Done():
				return
			case line, ok := <-lineChan:
				if !ok {
					return
				}
				fmt.Fprintf(w, "data: %s\n\n", string(line))
				flusher.Flush()
			}
		}
	}
}
