package utils

import "strings"

// SafeString mengembalikan string kosong jika nil
func SafeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// JoinAddress menggabungkan field alamat tanpa koma berlebih
func JoinAddress(parts ...*string) string {
	var result []string
	for _, p := range parts {
		if p != nil && *p != "" {
			result = append(result, *p)
		}
	}
	return strings.Join(result, ", ")
}
