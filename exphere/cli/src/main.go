package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/jonatassales/tao/exphere/cli/src/infrastructure/process"
	"github.com/jonatassales/tao/exphere/cli/src/presentation/cli"
	"github.com/jonatassales/tao/exphere/cli/src/presentation/cli/support"
)

// version can be injected at build time; an unversioned build stays explicit.
var version = "0.0.0-dev"

func main() {
	applicationContext, stopSignals := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)

	cases, _ := useCases()
	exitStatus := cli.Execute(
		applicationContext,
		os.Args[1:],
		support.Options{
			Version:     version,
			Input:       process.Input{Context: applicationContext, File: os.Stdin},
			Output:      os.Stdout,
			Diagnostics: os.Stderr,
		},
		cases,
	)

	stopSignals()
	os.Exit(exitStatus)
}
