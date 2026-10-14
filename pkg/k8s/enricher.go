package k8s

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/sujanv/LogSift/pkg/model"
)

type PodMetadata struct {
	PodName       string `json:"pod_name"`
	Namespace     string `json:"namespace"`
	ContainerName string `json:"container_name"`
	PodUID        string `json:"pod_uid,omitempty"`
}

func ExtractFromPath(logFilePath string) *PodMetadata {
	clean := filepath.ToSlash(logFilePath)
	parts := strings.Split(clean, "/")
	meta := &PodMetadata{}

	for i, part := range parts {
		if part == "pods" && i+2 < len(parts) {
			podPart := parts[i+1]
			meta.ContainerName = parts[i+2]
			sub := strings.Split(podPart, "_")
			if len(sub) >= 2 {
				meta.Namespace = sub[0]
				meta.PodName = sub[1]
			}
			if len(sub) >= 3 {
				meta.PodUID = sub[2]
			}
			return meta
		}
	}
	return nil
}

func ParseCRILine(line string, r *model.Record) bool {
	parts := strings.SplitN(line, " ", 4)
	if len(parts) < 4 {
		return false
	}
	if ts, err := time.Parse(time.RFC3339Nano, parts[0]); err == nil {
		r.Timestamp = ts
	}
	r.Message = parts[3]
	if r.Fields == nil {
		r.Fields = make(map[string]interface{})
	}
	r.Fields["k8s.stream"] = parts[1]
	r.Fields["k8s.log_tag"] = parts[2]
	return true
}
