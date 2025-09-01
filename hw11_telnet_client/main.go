package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	host string
	port string
)

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "timeout for connection")
	flag.Parse()

	if len(flag.Args()) >= 1 {
		host = flag.Args()[0]
	} else {
		host = "localhost"
	}

	if len(flag.Args()) >= 2 {
		port = flag.Args()[1]
	} else {
		port = "23"
	}

	address := net.JoinHostPort(host, port)
	client := NewTelnetClient(address, *timeout, os.Stdin, os.Stdout)
	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}

	go func() {
		err := client.Send()
		if err != nil {
			log.Fatal(err)
		}
	}()
	go func() {
		err := client.Receive()
		if err != nil {
			log.Fatal(err)
		}
	}()

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan bool, 1)

	go func() {
		defer client.Close()

		<-sigs
		done <- true
	}()

	<-done
}
