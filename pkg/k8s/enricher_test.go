package k8s

import (
	"testing"

	"github.com/sujanv/LogSift/pkg/model"
)

func TestExtractFromPath(t *testing.T) {
	path := "/var/log/pods/kube-system_coredns-789_uid-abc/coredns/0.log"
	meta := ExtractFromPath(path)
	if meta == nil {
		t.Fatalf("expected extracted pod metadata")
	}
	if meta.Namespace != "kube-system" || meta.PodName != "coredns-789" || meta.ContainerName != "coredns" {
		t.Errorf("mismatch: %+v", meta)
	}
}

func TestParseCRILine(t *testing.T) {
	line := "2026-10-14T10:00:00.123456789Z stdout F Starting HTTP server on port 8080"
	rec := &model.Record{}
	ok := ParseCRILine(line, rec)
	if !ok || rec.Message != "Starting HTTP server on port 8080" {
		t.Errorf("CRI line parse failed: %+v", rec)
	}
}
