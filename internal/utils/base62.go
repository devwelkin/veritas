package utils

import (
	"math"
)

const (
	base62Chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	base        = uint64(len(base62Chars))
)

// tobase62 converts a number to its base62 representation.
// this version is more efficient as it avoids the string reversal step.
func ToBase62(n uint64) string {
	if n == 0 {
		return string(base62Chars[0])
	}

	// calculate length of the result
	length := int(math.Floor(math.Log(float64(n))/math.Log(float64(base))) + 1)
	buf := make([]byte, length)

	// fill buffer from right to left
	i := length - 1
	for n > 0 {
		buf[i] = base62Chars[n%base]
		n /= base
		i--
	}

	return string(buf)
}
