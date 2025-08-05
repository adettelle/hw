package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	wantedEnv := Environment{
		"BAR":   EnvValue{Value: "bar", NeedRemove: false},
		"EMPTY": EnvValue{Value: "", NeedRemove: false},
		"FOO":   EnvValue{Value: "   foo\nwith new line", NeedRemove: false},
		"HELLO": EnvValue{Value: "\"hello\"", NeedRemove: false},
		"UNSET": EnvValue{Value: "", NeedRemove: true},
	}

	dir := "./testdata/env"

	env, err := ReadDir(dir)
	require.NoError(t, err)

	require.Equal(t, len(wantedEnv), len(env))

	for key := range wantedEnv {
		require.Equal(t, wantedEnv[key], env[key])
	}
}

func TestReadFirstLine(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		expected string
	}{
		{
			name:     "single line",
			data:     "sdsdf ",
			expected: "sdsdf ",
		},
		{
			name:     "middle enter",
			data:     "abc\ndef",
			expected: "abc",
		},
	}

	for _, tCase := range tests {
		t.Run(tCase.name, func(t *testing.T) {
			data := bytes.NewBuffer([]byte(tCase.data))
			result, err := readFirstLine(data)
			require.NoError(t, err)
			require.Equal(t, tCase.expected, result)
		})
	}
}

func TestNormalize(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		expected string
	}{
		{
			name:     "single space",
			data:     " ",
			expected: "",
		},
		{
			name:     "last tab",
			data:     " abc \t",
			expected: " abc",
		},
		{
			name:     "last space",
			data:     " abc  ",
			expected: " abc",
		},
		{
			name:     "zero byte",
			data:     "123\x00",
			expected: "123\n",
		},
	}

	for _, tCase := range tests {
		t.Run(tCase.name, func(t *testing.T) {
			res := normalize(tCase.data)
			require.Equal(t, tCase.expected, res)
		})
	}
}
