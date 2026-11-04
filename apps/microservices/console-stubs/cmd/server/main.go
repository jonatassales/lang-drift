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
	if err := run(ctx, configFromEnv(os.Getenv), os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "console-stubs:", err)
		os.Exit(1)
	}
}
