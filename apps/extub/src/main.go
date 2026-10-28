package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	if err := run(ctx, os.Getenv("EXTUB_PORT"), os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "extub:", err)
		os.Exit(1)
	}
}
