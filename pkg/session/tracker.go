package session

import (
	"sync"
	"time"

	"github.com/sujanv/LogSift/pkg/model"
)

type Session struct {
	ID        string        `json:"id"`
	ClientIP  string        `json:"client_ip"`
	StartTime time.Time     `json:"start_time"`
	EndTime   time.Time     `json:"end_time"`
	Duration  time.Duration `json:"duration"`
	Requests  int           `json:"requests"`
	Pages     []string      `json:"pages"`
}

type SessionTracker struct {
	mu        sync.Mutex
	timeout   time.Duration
	active    map[string]*Session
	completed []*Session
}

func NewSessionTracker(timeout time.Duration) *SessionTracker {
	if timeout <= 0 {
		timeout = 30 * time.Minute
	}
	return &SessionTracker{
		timeout: timeout,
		active:  make(map[string]*Session),
	}
}

func (st *SessionTracker) Ingest(r *model.Record) {
	key := r.IP
	if r.UserID != "" {
		key = r.UserID
	}
	if key == "" {
		return
	}

	st.mu.Lock()
	defer st.mu.Unlock()

	s, exists := st.active[key]
	if !exists || r.Timestamp.Sub(s.EndTime) > st.timeout {
		if exists {
			s.Duration = s.EndTime.Sub(s.StartTime)
			st.completed = append(st.completed, s)
		}
		s = &Session{
			ID:        key,
			ClientIP:  r.IP,
			StartTime: r.Timestamp,
			EndTime:   r.Timestamp,
			Requests:  1,
			Pages:     []string{r.Path},
		}
		st.active[key] = s
		return
	}

	s.Requests++
	s.EndTime = r.Timestamp
	s.Pages = append(s.Pages, r.Path)
}

func (st *SessionTracker) CompleteFinalizes() []*Session {
	st.mu.Lock()
	defer st.mu.Unlock()

	for _, s := range st.active {
		s.Duration = s.EndTime.Sub(s.StartTime)
		st.completed = append(st.completed, s)
	}
	st.active = make(map[string]*Session)
	return st.completed
}
