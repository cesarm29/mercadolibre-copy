package models

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

// StringSlice is a PostgreSQL text[] column that also marshals to a JSON array.
type StringSlice []string

// Scan implements sql.Scanner for reading a PG array literal like {"a","b"}.
func (s *StringSlice) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}
	var input string
	switch v := value.(type) {
	case string:
		input = v
	case []byte:
		input = string(v)
	default:
		return fmt.Errorf("StringSlice: cannot scan type %T", value)
	}
	out := []string{}
	p := strings.TrimSpace(input)
	p = strings.Trim(p, "{}")
	if p == "" {
		*s = out
		return nil
	}
	for _, item := range strings.Split(p, ",") {
		item = strings.TrimSpace(item)
		item = strings.Trim(item, "\"")
		out = append(out, item)
	}
	*s = out
	return nil
}

// Value implements driver.Valuer for writing a PG array literal.
func (s StringSlice) Value() (driver.Value, error) {
	parts := make([]string, 0, len(s))
	for _, v := range s {
		escaped := strings.ReplaceAll(v, `"`, `\"`)
		parts = append(parts, fmt.Sprintf("%q", escaped))
	}
	return "{" + strings.Join(parts, ",") + "}", nil
}