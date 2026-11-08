package embed

import (
	"testing"
)

func TestVectorizerAndCosine(t *testing.T) {
	v := NewVectorizer()
	corpus := []string{
		"connection reset by peer on socket 10.0.0.1",
		"failed to connect to postgresql database host",
		"disk quota exceeded on volume /var/log",
	}
	v.Train(corpus)

	v1 := v.Transform("connection reset by peer")
	v2 := v.Transform("connection reset on socket")
	v3 := v.Transform("unrelated disk space full error")

	sim12 := CosineSimilarity(v1, v2)
	sim13 := CosineSimilarity(v1, v3)

	if sim12 <= sim13 {
		t.Errorf("expected similar logs to have higher similarity: sim12=%f, sim13=%f", sim12, sim13)
	}
	if sim12 < 0.3 {
		t.Errorf("sim12 should be high, got %f", sim12)
	}
}
