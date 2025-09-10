package configs

import (
	"context"
	"flag"
	"os"
	"testing"

	"github.com/c2fo/testify/require"
)

func TestLoadConfig(t *testing.T) {
	ctx := context.Background()
	cfg, err := New(&ctx, true, "./cfg.json")
	require.NoError(t, err)

	require.Equal(t, "ERROR", cfg.Logger.Level)
	require.Equal(t, "localhost:8080", cfg.Address)
	require.Equal(t, "localhost", cfg.DBHost)
	require.Equal(t, "9999", cfg.DBPort)
	require.Equal(t, "postgres", cfg.DBUser)
	require.Equal(t, "123456", cfg.DBPassword)
	require.Equal(t, "calendar", cfg.DBName)
}

func TestLoadInexistentConfig(t *testing.T) {
	ctx := context.Background()
	_, err := New(&ctx, true, "./cfg1.json")
	require.Error(t, err)
}

func TestFromFlagsConfig(t *testing.T) {
	// сбрасываем глобальные флаги
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	os.Args = []string{"cmd", "-a=localhost:1234", "-h=localhost"}

	ctx := context.Background()

	cfg, err := New(&ctx, false, "")
	require.NoError(t, err)

	require.Equal(t, "INFO", cfg.Logger.Level)
	require.Equal(t, "localhost:1234", cfg.Address)
	require.Equal(t, "localhost", cfg.DBHost)
	require.Equal(t, "9999", cfg.DBPort)
	require.Equal(t, "postgres", cfg.DBUser)
	require.Equal(t, "123456", cfg.DBPassword)
	require.Equal(t, "calendar", cfg.DBName)
}
