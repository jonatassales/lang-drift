package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jonatassales/tao/apps/extub/src/application"
	"github.com/jonatassales/tao/apps/extub/src/infrastructure/probe"
	"github.com/jonatassales/tao/apps/extub/src/infrastructure/rtc"
	"github.com/jonatassales/tao/apps/extub/src/infrastructure/seed"
	"github.com/jonatassales/tao/apps/extub/src/presentation/httpapi"
)

// config is what the environment decides: the token that switches extub
// to its hosted preview, and the port it listens on.
type config struct {
	token     string
	port      string
	extubPort string
}

func configFromEnv(lookup func(string) string) config {
	return config{token: strings.TrimSpace(lookup("EXTUB_TOKEN")), port: lookup("PORT"), extubPort: lookup("EXTUB_PORT")}
}

// A hosted token shorter than this is refused at start, so a pasted
// fragment never becomes the only gate on the preview.
const minimumTokenLength = 32

// listenAddress is loopback on EXTUB_PORT, 3401 by default, for the local
// contract; with a token it is every interface on PORT, 3000 by default, as
// Vercel's Go preset requires.
func listenAddress(settings config) (string, error) {
	if settings.token != "" && len(settings.token) < minimumTokenLength {
		return "", fmt.Errorf("EXTUB_TOKEN must be at least %d characters", minimumTokenLength)
	}
	host, port, name := "127.0.0.1", settings.extubPort, "EXTUB_PORT"
	if port == "" {
		port = "3401"
	}
	if settings.token != "" {
		host, port, name = "", settings.port, "PORT"
		if port == "" {
			port = "3000"
		}
	}
	number, err := strconv.Atoi(port)
	if err != nil || number < 1 || number > 65535 {
		return "", fmt.Errorf("%s must be between 1 and 65535", name)
	}
	return net.JoinHostPort(host, port), nil
}

// run is the composition root and owns HTTP, the stream loop, and every peer.
// The hosted preview runs no stream and no peer.
func run(ctx context.Context, settings config, log io.Writer) error {
	address, err := listenAddress(settings)
	if err != nil {
		return err
	}
	_, port, _ := net.SplitHostPort(address)
	hosted := settings.token != ""
	store, err := seed.New(time.Now())
	if err != nil {
		return err
	}
	service := application.Service{Store: store, Now: time.Now}
	handler := httpapi.Handler{Token: settings.token}
	var hub *rtc.Hub
	if !hosted {
		service.LoadTests = &application.LoadTester{Probe: probe.New(port)}
		hub = rtc.New()
		defer hub.Close()
		handler.Signaler = hub
	}
	handler.Service = service
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer listener.Close()
	child, cancel := context.WithCancel(ctx)
	defer cancel()
	server := &http.Server{Handler: handler, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second, WriteTimeout: 15 * time.Second, IdleTimeout: 30 * time.Second, MaxHeaderBytes: 8192, BaseContext: func(net.Listener) context.Context { return child }}
	served := make(chan error, 1)
	streamed := make(chan error, 1)
	go func() { served <- server.Serve(listener) }()
	if hosted {
		// No stream runs; the channel settles when the server stops.
		go func() {
			<-child.Done()
			streamed <- nil
		}()
		fmt.Fprintf(log, "extub: hosted read-only preview on %s\n", address)
	} else {
		go func() { streamed <- service.Stream(child, hub) }()
		fmt.Fprintf(log, "extub: synthetic seeds only at http://%s\n", address)
	}
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
	if hub == nil {
		return result
	}
	return errors.Join(result, hub.Close())
}
