package utils

import (
	"strings"
)

const (
	base62Chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	base        = uint64(len(base62Chars))
)

// ToBase62 converts a number to its Base62 representation.
func ToBase62(n uint64) string {
	if n == 0 {
		return string(base62Chars[0])
	}

	var sb strings.Builder
	for n > 0 {
		sb.WriteByte(base62Chars[n%base])
		n /= base
	}

	// Reverse the string
	runes := []rune(sb.String())
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}
