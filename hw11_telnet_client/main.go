package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
)

func main() {
	timeout := flag.Duration("timeout", 0, "")
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		log.Fatal("host or port is not provided")
	}

	client := NewTelnetClient(net.JoinHostPort(args[0], args[1]), *timeout, os.Stdin, os.Stdout)
	if err := client.Connect(); err != nil {
		log.Fatalf("connection failure: %s", err)
	}
	defer client.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	go func() {
		if err := client.Send(); err != nil {
			fmt.Fprintf(os.Stderr, "command send error: %s", err)
		} else {
			fmt.Fprint(os.Stderr, "..EOF")
		}
		cancel()
	}()

	go func() {
		if err := client.Receive(); err != nil {
			fmt.Fprintf(os.Stderr, "command receive error: %s", err)
		} else {
			fmt.Fprint(os.Stderr, "..Connection was closed by peer")
		}
		cancel()
	}()

	<-ctx.Done()
}
