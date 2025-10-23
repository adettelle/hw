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
	cfg, err := New(&ctx)
	require.NoError(t, err)

	require.Equal(t, "INFO", cfg.Logger.Level)
	require.Equal(t, "localhost:8081", cfg.Address)
	require.Equal(t, "localhost", cfg.DBHost)
	require.Equal(t, "9999", cfg.DBPort)
	require.Equal(t, "postgres", cfg.DBUser)
	require.Equal(t, "123456", cfg.DBPassword)
	require.Equal(t, "test_db", cfg.DBName)
}

func TestFromFlagsConfig(t *testing.T) {
	// сбрасываем глобальные флаги
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	os.Args = []string{"cmd", "-config=./calendar_cfg_test.json"}

	ctx := context.Background()

	cfg, err := New(&ctx)
	require.NoError(t, err)

	require.Equal(t, "INFO", cfg.Logger.Level)
	require.Equal(t, "localhost:8888", cfg.Address)
	require.Equal(t, "localhost", cfg.DBHost)
	require.Equal(t, "9999", cfg.DBPort)
	require.Equal(t, "postgres", cfg.DBUser)
	require.Equal(t, "123456", cfg.DBPassword)
	require.Equal(t, "test_db", cfg.DBName)
}
