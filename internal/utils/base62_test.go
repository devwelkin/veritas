package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToBase62(t *testing.T) {
	testCases := []struct {
		name     string
		input    uint64
		expected string
	}{
		{
			name:     "Test with zero",
			input:    0,
			expected: "a",
		},
		{
			name:     "Test with a single digit in base62",
			input:    10,
			expected: "k",
		},
		{
			name:     "Test with a number resulting in a two-character string",
			input:    62,
			expected: "ba",
		},
		{
			name:     "Test with a larger number",
			input:    12345,
			expected: "dnh",
		},
		{
			name:     "Test with another large number",
			input:    7891234,
			expected: "Hg18",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, ToBase62(tc.input))
		})
	}
}
