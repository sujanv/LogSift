package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/sujanv/LogSift/pkg/stats"
)

// Server provides an embedded monitoring HTTP API for LogSift.
type Server struct {
	addr       string
	mu         sync.RWMutex
	latTracker *stats.LatencyTracker
	errSum     *stats.ErrorSummarizer
	totalLogs  int64
	startTime  time.Time
}

// NewServer creates an HTTP monitoring server.
func NewServer(addr string, lat *stats.LatencyTracker, errs *stats.ErrorSummarizer) *Server {
	return &Server{
		addr:       addr,
		latTracker: lat,
		errSum:     errs,
		startTime:  time.Now(),
	}
}

// Routes registers the REST endpoints.
func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/api/v1/metrics", s.handleMetrics)
	mux.HandleFunc("/api/v1/errors", s.handleErrors)
	return mux
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func (s *Server) handleMetrics(w http.ResponseWriter, r *http.Request) {
	snap := s.latTracker.Snapshot()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(snap)
}

func (s *Server) handleErrors(w http.ResponseWriter, r *http.Request) {
	clusters := s.errSum.TopClusters(10)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(clusters)
}

// Start begins listening.
func (s *Server) Start() error {
	fmt.Printf("LogSift HTTP Server listening on %s\n", s.addr)
	return http.ListenAndServe(s.addr, s.Routes())
}
