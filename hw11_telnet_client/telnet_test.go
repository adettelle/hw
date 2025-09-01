package main

import (
	"bytes"
	"context"
	"io"
	"net"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTelnetClient(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		lc := net.ListenConfig{}
		l, err := lc.Listen(ctx, "tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, io.NopCloser(in), out)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			in.WriteString("hello\n")
			err = client.Send()
			require.NoError(t, err)

			err = client.Receive()
			require.NoError(t, err)
			require.Equal(t, "world\n", out.String())
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)
			require.Equal(t, "hello\n", string(request)[:n])

			n, err = conn.Write([]byte("world\n"))
			require.NoError(t, err)
			require.NotEqual(t, 0, n)
		}()

		wg.Wait()
	})
}

func TestConnectError(t *testing.T) {
	addr := "127.0.0.1:65000"
	client := NewTelnetClient(addr, 200*time.Millisecond,
		io.NopCloser(strings.NewReader("")), &bytes.Buffer{})

	err := client.Connect()
	require.Error(t, err, "expected an error when connecting to a non-existent server")
}

func TestSendWithoutConnect(t *testing.T) {
	client := &TClient{}
	err := client.Send()
	require.Error(t, err)
	require.Contains(t, err.Error(), "not connected")
}

func TestReceiveWithoutConnect(t *testing.T) {
	client := &TClient{}
	err := client.Receive()
	require.Error(t, err)
	require.Contains(t, err.Error(), "not connected")
}

func TestCloseWithoutConnect(t *testing.T) {
	client := &TClient{}
	err := client.Close()
	require.NoError(t, err)
}
