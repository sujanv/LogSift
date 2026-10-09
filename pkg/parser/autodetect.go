package parser

import (
	"bytes"

	"github.com/sujanv/LogSift/pkg/model"
)

// AutoDetect examines a sample of lines to determine the most likely log parser format.
func AutoDetect(sampleLines [][]byte) Parser {
	if len(sampleLines) == 0 {
		return NewJSONParser()
	}

	jsonParser := NewJSONParser()
	logfmtParser := NewLogfmtParser()
	clfParser := NewCLFParser()
	syslogParser := NewSyslogParser()

	rec := model.DefaultRecordPool.Get()
	defer model.DefaultRecordPool.Put(rec)

	jsonScore := 0
	logfmtScore := 0
	clfScore := 0
	syslogScore := 0

	for _, line := range sampleLines {
		trimmed := bytes.TrimSpace(line)
		if len(trimmed) == 0 {
			continue
		}

		// Quick JSON heuristic
		if trimmed[0] == '{' && trimmed[len(trimmed)-1] == '}' {
			rec.Reset()
			if err := jsonParser.Parse(trimmed, rec); err == nil && len(rec.Fields) > 0 {
				jsonScore += 5
			}
		}

		// Quick CLF heuristic
		if bytes.Contains(trimmed, []byte("] \"")) {
			rec.Reset()
			if err := clfParser.Parse(trimmed, rec); err == nil && rec.Method != "" {
				clfScore += 5
			}
		}

		// Syslog heuristic
		if trimmed[0] == '<' && bytes.Contains(trimmed, []byte(">")) {
			rec.Reset()
			if err := syslogParser.Parse(trimmed, rec); err == nil {
				syslogScore += 4
			}
		}

		// Logfmt heuristic: multiple key=val pairs
		if bytes.Count(trimmed, []byte("=")) >= 2 {
			rec.Reset()
			if err := logfmtParser.Parse(trimmed, rec); err == nil && len(rec.Fields) >= 2 {
				logfmtScore += 3
			}
		}
	}

	if jsonScore > 0 && jsonScore >= logfmtScore && jsonScore >= clfScore && jsonScore >= syslogScore {
		return jsonParser
	}
	if clfScore > 0 && clfScore >= logfmtScore && clfScore >= syslogScore {
		return clfParser
	}
	if syslogScore > 0 && syslogScore >= logfmtScore {
		return syslogParser
	}
	if logfmtScore > 0 {
		return logfmtParser
	}

	return jsonParser
}
