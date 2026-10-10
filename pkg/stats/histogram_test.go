package stats

import (
	"strings"
	"testing"
	"time"
)

func TestTimeHistogram(t *testing.T) {
	th := NewTimeHistogram(time.Minute)
	t0 := time.Date(2026, 10, 9, 12, 0, 10, 0, time.UTC)

	th.Add(t0, false)
	th.Add(t0.Add(10*time.Second), true)
	th.Add(t0.Add(2*time.Minute), false)

	buckets := th.Buckets()
	if len(buckets) != 2 {
		t.Fatalf("expected 2 buckets, got %d", len(buckets))
	}

	if buckets[0].Total != 2 || buckets[0].Errors != 1 {
		t.Errorf("bucket 0: total=%d, errors=%d", buckets[0].Total, buckets[0].Errors)
	}

	ascii := th.RenderASCII(40)
	if !strings.Contains(ascii, "reqs") {
		t.Errorf("expected ascii chart output, got:\n%s", ascii)
	}
}
