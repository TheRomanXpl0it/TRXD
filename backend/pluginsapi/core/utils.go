package core

import (
	"fmt"
	"reflect"
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