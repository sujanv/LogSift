package generator

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"time"
)

var (
	endpoints = []string{
		"/api/v1/auth/login",
		"/api/v1/users/profile",
		"/api/v1/items/search",
		"/api/v1/orders/checkout",
		"/api/v1/payments/process",
		"/healthz",
	}
	methods = []string{"GET", "POST", "PUT", "DELETE"}
	errorMsgs = []string{
		"connection reset by peer during TLS handshake with db-%d.internal",
		"failed to acquire connection pool lock within timeout %dms",
		"rate limit exceeded for client token usr-%05d",
		"downstream payment gateway responded with HTTP 504 Gateway Timeout",
		"invalid user session token: expired at %s",
	}
)

// Generator produces high-speed synthetic logs for benchmarking and load testing.
type Generator struct {
	rnd *rand.Rand
}

// NewGenerator creates a log generator.
func NewGenerator(seed int64) *Generator {
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	return &Generator{rnd: rand.New(rand.NewSource(seed))}
}

// GenerateJSONLines writes N synthetic JSON log lines into w.
func (g *Generator) GenerateJSONLines(w io.Writer, count int64) error {
	baseTime := time.Now().Add(-time.Hour)

	for i := int64(0); i < count; i++ {
		ts := baseTime.Add(time.Duration(i) * 5 * time.Millisecond)
		ep := endpoints[g.rnd.Intn(len(endpoints))]
		method := methods[g.rnd.Intn(len(methods))]

		// Log-normal latency distribution
		latMs := 10.0 + g.rnd.ExpFloat64()*45.0
		isError := g.rnd.Float64() < 0.05 // 5% error rate

		level := "info"
		status := 200
		msg := fmt.Sprintf("%s %s completed successfully", method, ep)

		if isError {
			level = "error"
			status = 500 + g.rnd.Intn(5)
			template := errorMsgs[g.rnd.Intn(len(errorMsgs))]
			msg = fmt.Sprintf(template, g.rnd.Intn(100), ts.Format(time.Kitchen))
			latMs += 500.0 // latency spike on error
		}

		payload := map[string]interface{}{
			"timestamp":  ts.Format(time.RFC3339Nano),
			"level":      level,
			"msg":        msg,
			"path":       ep,
			"method":     method,
			"status":     status,
			"latency_ms": latMs,
			"client_ip":  fmt.Sprintf("192.168.1.%d", 10+g.rnd.Intn(200)),
			"trace_id":   fmt.Sprintf("trace-%x-%04d", ts.Unix(), i%10000),
		}

		data, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		data = append(data, '\n')
		if _, err := w.Write(data); err != nil {
			return err
		}
	}

	return nil
}
