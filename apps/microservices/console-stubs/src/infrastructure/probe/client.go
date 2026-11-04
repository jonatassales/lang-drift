// Package probe measures requests to the fixed local stub, never user-supplied targets.
package probe

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

type Client struct {
	url    string
	client *http.Client
}

// New binds the probe to the same loopback listener as extub. It disables proxies,
// redirects and retries, and caps sockets and response headers for the 50 workers.
func New(port string) *Client {
	transport := &http.Transport{
		Proxy:                  nil,
		DialContext:            (&net.Dialer{Timeout: time.Second}).DialContext,
		DisableKeepAlives:      true,
		MaxConnsPerHost:        50,
		MaxResponseHeaderBytes: 1024,
		ResponseHeaderTimeout:  time.Second,
	}
	return &Client{
		url:    "http://" + net.JoinHostPort("127.0.0.1", port) + "/probe",
		client: &http.Client{Transport: transport, Timeout: time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }},
	}
}

func (client *Client) Check(ctx context.Context) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, client.url, nil)
	if err != nil {
		return err
	}
	response, err := client.client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("stub probe returned status %d", response.StatusCode)
	}
	return nil
}
