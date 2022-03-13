package convert

import (
	"fmt"
	"strconv"
)

func RawValueString(v interface{}) string {
	// Check if value is <nil>
	if v == nil {
		return ""
	}
	// Check other types
	switch out := v.(type) {
	case string:
		return out
	case int:
		return strconv.Itoa(out)
	case float64:
		return strconv.FormatFloat(out, 'e', -1, 64)
	case bool:
		return strconv.FormatBool(out)
	default:
		return fmt.Sprintf("%v", out)
	}
}
