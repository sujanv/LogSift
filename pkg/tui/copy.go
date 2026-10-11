package tui

import (
	"encoding/json"

	"github.com/atotto/clipboard"
	"github.com/sujanv/LogSift/pkg/model"
)

// CopyRecordToClipboard formats and writes record to system clipboard.
func CopyRecordToClipboard(r *model.Record) error {
	if r == nil {
		return nil
	}
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return clipboard.WriteAll(r.Message)
	}
	return clipboard.WriteAll(string(data))
}
