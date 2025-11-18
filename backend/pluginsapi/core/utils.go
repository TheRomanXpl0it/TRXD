package core

import (
	"fmt"
	"reflect"
	"strings"
)

func CheckSameTypes(before, after any) error {

	beforetype := reflect.TypeOf(before)
	aftertype := reflect.TypeOf(after)

	if beforetype == nil || aftertype == nil {
		if beforetype != aftertype {
			return fmt.Errorf("type changed: %v -> %v (nil mismatch)", beforetype, aftertype)
		}
	}
	if beforetype != aftertype {
		return fmt.Errorf("type changed: %v -> %v", beforetype, aftertype)
	}
	
	return nil
}

func ExtractSqlcName(query string) string {
	lines := strings.SplitN(query, "\n", 2)
	if len(lines) == 0 {
		return ""
	}
	line := strings.TrimSpace(lines[0])
	if !strings.HasPrefix(line, "-- name:") {
		return ""
	}

	rest := strings.TrimSpace(strings.TrimPrefix(line, "-- name:"))
	parts := strings.Fields(rest)
	if len(parts) == 0 {
		return ""
	}
	return parts[0]
}