package utils

import (
	"reflect"
	"testing"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "Empty string",
			input:    "",
			expected: []string{},
		},
		{
			name:     "Single argument",
			input:    "foo",
			expected: []string{"foo"},
		},
		{
			name:     "Multiple arguments",
			input:    "foo bar baz",
			expected: []string{"foo", "bar", "baz"},
		},
		{
			name:     "Arguments with quoted string",
			input:    `foo "bar baz"`,
			expected: []string{"foo", "bar baz"},
		},
		{
			name:     "Multiple quoted strings",
			input:    `"hello world" "foo bar"`,
			expected: []string{"hello world", "foo bar"},
		},
		{
			name:     "Mixed quoted and unquoted",
			input:    `ls -la "/home/user/my documents" test.txt`,
			expected: []string{"ls", "-la", "/home/user/my documents", "test.txt"},
		},
		{
			name:     "Leading and trailing spaces",
			input:    "  foo   bar  ",
			expected: []string{"foo", "bar"},
		},
		{
			name:     "Only spaces",
			input:    "   ",
			expected: []string{},
		},
		{
			name:     "Quoted string with special characters",
			input:    `"hello@world.com" "test-file.txt"`,
			expected: []string{"hello@world.com", "test-file.txt"},
		},
		{
			name:     "Command with flags and quoted arguments",
			input:    `grep -r "search term" /path/to/dir`,
			expected: []string{"grep", "-r", "search term", "/path/to/dir"},
		},
		{
			name:     "Empty quotes",
			input:    `foo "" bar`,
			expected: []string{"foo", "bar"}, // Empty strings between quotes are not added
		},
		{
			name:     "Single quoted word",
			input:    `"word"`,
			expected: []string{"word"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseArgs(tt.input)

			// Handle nil vs empty slice comparison
			if len(result) == 0 && len(tt.expected) == 0 {
				return // Both are effectively empty
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ParseArgs(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func BenchmarkParseArgs(b *testing.B) {
	input := `command -flag1 "argument with spaces" -flag2=value "another argument"`
	for i := 0; i < b.N; i++ {
		ParseArgs(input)
	}
}
