package slowlog

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

type SlowQuery struct {
	Database      string        `json:"database"`
	Timestamp     time.Time     `json:"timestamp"`
	Duration      time.Duration `json:"duration"`
	Query         string        `json:"query"`
	NormalizedSQL string        `json:"normalized_sql"`
	User          string        `json:"user"`
	Host          string        `json:"host"`
	RowsExamined  int64         `json:"rows_examined"`
}

var (
	pgRegex         = regexp.MustCompile(`duration:\s*([\d\.]+)\s*ms\s+statement:\s*(.*)`)
	mysqlTimeRegex  = regexp.MustCompile(`# Query_time:\s*([\d\.]+)\s+Lock_time:\s*([\d\.]+)\s+Rows_sent:\s*(\d+)\s+Rows_examined:\s*(\d+)`)
	numRegex        = regexp.MustCompile(`\b\d+\b`)
	strLiteralRegex = regexp.MustCompile(`'[^']*'`)
)

func ParsePostgresLine(line string) *SlowQuery {
	matches := pgRegex.FindStringSubmatch(line)
	if len(matches) < 3 {
		return nil
	}
	ms, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return nil
	}
	q := strings.TrimSpace(matches[2])
	return &SlowQuery{
		Database:      "postgresql",
		Duration:      time.Duration(ms * float64(time.Millisecond)),
		Query:         q,
		NormalizedSQL: NormalizeSQL(q),
	}
}

func ParseMySQLHeader(headerLine, queryLine string) *SlowQuery {
	matches := mysqlTimeRegex.FindStringSubmatch(headerLine)
	if len(matches) < 5 {
		return nil
	}
	sec, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return nil
	}
	rowsExamined, _ := strconv.ParseInt(matches[4], 10, 64)
	q := strings.TrimSpace(queryLine)
	return &SlowQuery{
		Database:      "mysql",
		Duration:      time.Duration(sec * float64(time.Second)),
		Query:         q,
		NormalizedSQL: NormalizeSQL(q),
		RowsExamined:  rowsExamined,
	}
}

func NormalizeSQL(sql string) string {
	sql = strLiteralRegex.ReplaceAllString(sql, "?")
	sql = numRegex.ReplaceAllString(sql, "?")
	sql = strings.Join(strings.Fields(sql), " ")
	return strings.ToUpper(sql)
}
