package reader

import (
	"os"
	"path/filepath"
	"strings"
)

// DiscoverFiles expands glob patterns and searches directories recursively for log files.
func DiscoverFiles(patterns []string) ([]string, error) {
	var results []string
	seen := make(map[string]bool)

	for _, pat := range patterns {
		// Check direct directory
		info, err := os.Stat(pat)
		if err == nil && info.IsDir() {
			err = filepath.Walk(pat, func(p string, f os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !f.IsDir() && isLogFile(p) && !seen[p] {
					seen[p] = true
					results = append(results, p)
				}
				return nil
			})
			if err != nil {
				return nil, err
			}
			continue
		}

		// Glob pattern
		matches, err := filepath.Glob(pat)
		if err != nil {
			return nil, err
		}
		for _, m := range matches {
			if !seen[m] {
				seen[m] = true
				results = append(results, m)
			}
		}
	}

	return results, nil
}

func isLogFile(p string) bool {
	lower := strings.ToLower(p)
	return strings.HasSuffix(lower, ".log") ||
		strings.HasSuffix(lower, ".log.gz") ||
		strings.HasSuffix(lower, ".json") ||
		strings.HasSuffix(lower, ".txt")
}
