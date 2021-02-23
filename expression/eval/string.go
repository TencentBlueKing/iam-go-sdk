package eval

import "strings"

func convertArgsToString(v1, v2 interface{}) (s1, s2 string, ok bool) {
	s1, ok = v1.(string)
	if !ok {
		return "", "", false
	}

	s2, ok = v2.(string)
	if !ok {
		return "", "", false
	}

	return s1, s2, true
}

// StartsWith return true if v1 startswith v2
func StartsWith(v1 interface{}, v2 interface{}) bool {
	s1, s2, ok := convertArgsToString(v1, v2)
	if !ok {
		return false
	}

	return strings.HasPrefix(s1, s2)
}

// NotStartsWith return true if v1 not startswith v2
func NotStartsWith(v1 interface{}, v2 interface{}) bool {
	s1, s2, ok := convertArgsToString(v1, v2)
	if !ok {
		return false
	}

	return !strings.HasPrefix(s1, s2)
}

// EndsWith return true if v1 endswith v2
func EndsWith(v1 interface{}, v2 interface{}) bool {
	s1, s2, ok := convertArgsToString(v1, v2)
	if !ok {
		return false
	}

	return strings.HasSuffix(s1, s2)
}

// NotEndsWith return true if v1 not endswith v2
func NotEndsWith(v1 interface{}, v2 interface{}) bool {
	s1, s2, ok := convertArgsToString(v1, v2)
	if !ok {
		return false
	}

	return !strings.HasSuffix(s1, s2)
}
