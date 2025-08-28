package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type TClient struct {
	Network string
	Timeout time.Duration
	Address string
	conn    net.Conn
	In      io.ReadCloser
	Out     io.Writer
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &TClient{Network: "tcp", Timeout: timeout, Address: address, In: in, Out: out}
}

func (tc *TClient) Connect() error {
	if tc.conn != nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), tc.Timeout)
	defer cancel()

	dialer := net.Dialer{}
	conn, err := dialer.DialContext(ctx, tc.Network, tc.Address)
	if err != nil {
		fmt.Fprintf(os.Stderr, "connection to %s failed: %v\n", tc.Address, err)
		return err
	}
	tc.conn = conn

	return nil
}

func (tc *TClient) Close() error {
	if tc.conn == nil {
		return nil
	}
	err := tc.conn.Close()
	tc.conn = nil
	return err
}

func (tc *TClient) Send() error {
	if tc.conn == nil {
		return fmt.Errorf("not connected")
	}
	// stdin -> socket
	_, err := io.Copy(tc.conn, tc.In)
	return err
}

func (tc *TClient) Receive() error {
	if tc.conn == nil {
		return fmt.Errorf("not connected")
	}
	// socket -> stdout
	_, err := io.Copy(tc.Out, tc.conn)
	return err
}
