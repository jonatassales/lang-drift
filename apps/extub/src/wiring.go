package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/jonatassales/tao/apps/extub/src/application"
	"github.com/jonatassales/tao/apps/extub/src/infrastructure/probe"
	"github.com/jonatassales/tao/apps/extub/src/infrastructure/rtc"
	"github.com/jonatassales/tao/apps/extub/src/infrastructure/seed"
	"github.com/jonatassales/tao/apps/extub/src/presentation/httpapi"
)

// run is the composition root and owns HTTP, the stream loop, and every peer.
func run(ctx context.Context, port string, log io.Writer) error {
	if port == "" {
		port = "3401"
	}
	number, err := strconv.Atoi(port)
	if err != nil || number < 1 || number > 65535 {
		return errors.New("EXTUB_PORT must be between 1 and 65535")
	}
	store, err := seed.New(time.Now())
	if err != nil {
		return err
	}
	service := application.Service{Store: store, Now: time.Now, LoadTests: &application.LoadTester{Probe: probe.New(port)}}
	hub := rtc.New()
	defer hub.Close()
	listener, err := net.Listen("tcp", net.JoinHostPort("127.0.0.1", port))
	if err != nil {
		return err
	}
	defer listener.Close()
	child, cancel := context.WithCancel(ctx)
	defer cancel()
	server := &http.Server{Handler: httpapi.Handler{Service: service, Signaler: hub}, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second, WriteTimeout: 15 * time.Second, IdleTimeout: 30 * time.Second, MaxHeaderBytes: 8192, BaseContext: func(net.Listener) context.Context { return child }}
	served := make(chan error, 1)
	streamed := make(chan error, 1)
	go func() { served <- server.Serve(listener) }()
	go func() { streamed <- service.Stream(child, hub) }()
	fmt.Fprintf(log, "extub: synthetic seeds only at http://127.0.0.1:%s\n", port)
	var result error
	httpDone, streamDone := false, false
	select {
	case <-ctx.Done():
	case result = <-served:
		httpDone = true
	case result = <-streamed:
		streamDone = true
	}
	cancel()
	shutdown, cancelShutdown := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancelShutdown()
	if err = server.Shutdown(shutdown); err != nil {
		result = errors.Join(result, err, server.Close())
	}
	if !httpDone {
		err = <-served
		if !errors.Is(err, http.ErrServerClosed) {
			result = errors.Join(result, err)
		}
	}
	if !streamDone {
		result = errors.Join(result, <-streamed)
	}
	return errors.Join(result, hub.Close())
}
