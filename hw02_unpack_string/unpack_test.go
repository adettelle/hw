package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},
		{input: "🙃0", expected: ""},
		{input: "aaф0b", expected: "aab"},

		// uncomment if task with asterisk completed
		{input: `qwe\4\5`, expected: `qwe45`},
		{input: `qwe\45`, expected: `qwe44444`},
		{input: `qwe\\5`, expected: `qwe\\\\\`},
		{input: `qwe\\\3`, expected: `qwe\3`},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b"}
	for _, tc := range invalidStrings {
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}

func TestUnpack2(t *testing.T) {
	tests := []struct {
		s, expected string
		err         error
	}{
		{"abcd", "abcd", nil},
		{"d\n5abc", "d\n\n\n\n\nabc", nil},
		{`qw\ne`, "", ErrInvalidString},

		{"ab\\72", "ab77", nil},
		{"ab\\x3", "", ErrInvalidString},
		{"d\abc", "d\abc", nil},
		{`qwe\4\`, "", ErrInvalidString},
		{"⌘こんにちは", "⌘こんにちは", nil},
		{"⌘2こんにちは", "⌘⌘こんにちは", nil},
		{"⌘0こんにちは", "こんにちは", nil},
		{"привет", "привет", nil},
		{"0привет", "", ErrInvalidString},
		{"п2р3ивет", "ппррривет", nil},
		{"привет0", "приве", nil},
		{"Привет, мир!", "Привет, мир!", nil},
		{"d\\2abc", `d\2abc`, nil},
	}

	for _, tt := range tests {
		result, err := Unpack(tt.s)

		if result != tt.expected && !errors.Is(err, tt.err) {
			t.Errorf("result: %s, want: %s\n", result, tt.expected)
		}
	}
}
