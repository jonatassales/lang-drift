//go:build e2e

package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/jonatassales/tao/packages/ex/src/infrastructure/process"
	"github.com/jonatassales/tao/packages/ex/src/presentation/cli/support"
)

// TestNativeLifecycleFixture is entered only by the opt-in E2E harness's
// subprocess. It exercises retained native adapters, not retired public CLI
// commands. No production flag, environment variable or backdoor is added.
func TestNativeLifecycleFixture(t *testing.T) {
	if os.Getenv("LANGDRIFT_NATIVE_FIXTURE") != "1" {
		t.Skip("native fixture subprocess only")
	}
	args := os.Args
	for index, arg := range args {
		if arg == "--" {
			args = args[index+1:]
			break
		}
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	status := executeNativeCatalogs(ctx, args, support.Options{Version: version, Input: process.Input{Context: ctx, File: os.Stdin}, Output: os.Stdout, Diagnostics: os.Stderr}, useCases())
	cancel()
	os.Exit(status)
}
