package main

import (
	"bytes"
	"io"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	envs := Environment{
		"BAR":   EnvValue{Value: "bar", NeedRemove: false},
		"EMPTY": EnvValue{Value: "", NeedRemove: false},
		"FOO":   EnvValue{Value: "   foo", NeedRemove: false},
		"HELLO": EnvValue{Value: "\"hello\"", NeedRemove: false},
		"UNSET": EnvValue{Value: "", NeedRemove: true},
	}
	require.NoError(t, os.Setenv("ADDED", "new"))

	expectedOutput := `HELLO is ("hello")
BAR is (bar)
FOO is (   foo)
UNSET is ()
ADDED is (new)
EMPTY is ()
arguments are arg1=1 arg2=2
`
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() { os.Stdout = old }()

	shell := "/bin/bash"

	args := []string{shell, path.Join("./testdata", "echo.sh"), "arg1=1", "arg2=2"}
	resCode := RunCmd(args, envs)
	require.Equal(t, 0, resCode)

	require.NoError(t, w.Close())

	var buf bytes.Buffer
	_, err := io.Copy(&buf, r)
	require.NoError(t, err)

	require.Equal(t, expectedOutput, buf.String())
}

func TestRunCmdNegative(t *testing.T) {
	resCode := RunCmd([]string{}, Environment{})
	require.Equal(t, -1, resCode)
}
