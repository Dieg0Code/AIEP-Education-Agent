package date

import "time"

const (
	LayoutDate     = "2006-01-02"
	LayoutDateTime = time.RFC3339
)

// ParseDate parses a date-only string (YYYY-MM-DD) into time.Time.
func ParseDate(s string) (time.Time, error) {
	return time.Parse(LayoutDate, s)
}

// FormatDate formats a time.Time into a date-only string (YYYY-MM-DD).
func FormatDate(t time.Time) string {
	return t.Format(LayoutDate)
}

// ParseDateTime parses an RFC3339 datetime string.
func ParseDateTime(s string) (time.Time, error) {
	return time.Parse(LayoutDateTime, s)
}

// FormatDateTime formats time.Time into RFC3339.
func FormatDateTime(t time.Time) string {
	return t.Format(LayoutDateTime)
}
