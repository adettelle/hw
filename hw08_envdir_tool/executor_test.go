package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmdNegative(t *testing.T) {
	resCode := RunCmd([]string{}, Environment{})
	require.Equal(t, -1, resCode)
}
