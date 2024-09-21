package swagger

import (
	"unicode"
)

func firstLetterToLower(s string) string {
	if len(s) == 0 {
		return s
	}

	r := []rune(s)
	r[0] = unicode.ToLower(r[0])

	return string(r)
}

func IsNil(val interface{}) bool {
	switch v := val.(type) {
	case string:
		return len(v) == 0
	case []string:
		return len(v) == 0
	case []*interface{}:
		return len(v) == 0
	case []interface{}:
		return len(v) == 0
	case map[string]interface{}:
		return len(v) == 0
	case []*SecuritySchemeObject:
		return len(v) == 0
	case []*ParameterObject:
		return len(v) == 0
	default:
		return val == nil
	}
}
