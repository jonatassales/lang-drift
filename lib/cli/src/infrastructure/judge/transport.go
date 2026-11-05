package judge

import (
	"errors"
	"net"
	"net/http"
	"net/url"
)

// maxResponseBytes bounds what the judge reads of one response, as the
// connected agent bounds its own.
const maxResponseBytes = 4 << 20

var (
	// errInsecureURL refuses a console URL that would carry the secret in
	// clear text to another machine.
	errInsecureURL = errors.New("the console URL must be HTTPS, or HTTP on loopback")

	// errRedirect refuses a redirect, which would carry the secret to the
	// host it names.
	errRedirect = errors.New("the console answered with a redirect, and redirects are refused")
)

// httpClient carries the judge's calls. It never follows a redirect, which
// would carry the bearer secret to the host it names.
func httpClient() *http.Client {
	return &http.Client{
		Transport: transport{next: http.DefaultTransport},
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return errRedirect
		},
	}
}

// transport refuses a call to a URL that may not carry the secret before
// anything is sent, and bounds what can be read of every response.
type transport struct{ next http.RoundTripper }

func (transport transport) RoundTrip(request *http.Request) (*http.Response, error) {
	if !secure(request.URL) {
		return nil, errInsecureURL
	}

	response, err := transport.next.RoundTrip(request)
	if err != nil {
		return nil, err
	}

	response.Body = http.MaxBytesReader(nil, response.Body, maxResponseBytes)

	return response, nil
}

// secure reports whether a URL may carry the secret: HTTPS, or HTTP to a
// loopback host for local development.
func secure(endpoint *url.URL) bool {
	if endpoint.Scheme == "https" {
		return true
	}

	address := net.ParseIP(endpoint.Hostname())
	loopback := endpoint.Hostname() == "localhost" || (address != nil && address.IsLoopback())

	return endpoint.Scheme == "http" && loopback
}
