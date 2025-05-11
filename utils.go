package swagger

import (
	"reflect"
	"time"
)

// IsNil checks if a given value is nil or empty.
// The function returns true for empty strings, slices, maps, and pointers.
// For other types, it returns whether the value is nil or not.
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

// mappingType takes a reflect.Value and returns a string describing its type in
// OpenAPI mapping terms. The returned string is one of "boolean", "integer",
// "number", "string", or "object".
func mappingType(val reflect.Type) string {
	if val == reflect.TypeOf(time.Time{}) {
		return "string"
	}
	switch val.Kind() {
	case reflect.Bool:
		return "boolean"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "integer"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "integer"
	case reflect.Float32, reflect.Float64:
		return "number"
	case reflect.String:
		return "string"
	case reflect.Pointer:
		return "object"
	case reflect.Slice, reflect.Array:
		return "array"
	default:
		return ""
	}
}
