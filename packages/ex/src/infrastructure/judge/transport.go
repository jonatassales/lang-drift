package judge

import (
	"errors"
	"net"
	"net/http"
	"net/url"
	"os"
)

// maxResponseBytes bounds what the judge reads of one response, as the
// console client bounds its own.
const maxResponseBytes = 4 << 20

var (
	// errInsecureBaseURL refuses a base URL that would carry the key in clear
	// text to another machine.
	errInsecureBaseURL = errors.New("the model's base URL must be HTTPS, or HTTP on loopback")

	// errRedirect refuses a redirect, which would carry the key to the host
	// it names.
	errRedirect = errors.New("the API answered with a redirect, and redirects are refused")

	// errOversized names a response cut at maxResponseBytes.
	errOversized = errors.New("the API's response exceeds 4 MiB")
)

// httpClient carries the judge's calls. It never follows a redirect: the SDK
// sends the key as X-Api-Key, a header Go keeps when a redirect leads to
// another host.
func httpClient() *http.Client {
	return &http.Client{
		Transport: transport{next: http.DefaultTransport},
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return errRedirect
		},
	}
}

// transport refuses a call to a URL that may not carry the key before
// anything is sent, whichever source set the base URL, and bounds what can
// be read of every response.
type transport struct{ next http.RoundTripper }

func (transport transport) RoundTrip(request *http.Request) (*http.Response, error) {
	if !secure(request.URL) {
		return nil, errInsecureBaseURL
	}

	response, err := transport.next.RoundTrip(request)
	if err != nil {
		return nil, err
	}

	response.Body = http.MaxBytesReader(nil, response.Body, maxResponseBytes)

	return response, nil
}

// secure reports whether a URL may carry the key: HTTPS, or HTTP to a
// loopback host for local development.
func secure(endpoint *url.URL) bool {
	if endpoint.Scheme == "https" {
		return true
	}

	address := net.ParseIP(endpoint.Hostname())
	loopback := endpoint.Hostname() == "localhost" || (address != nil && address.IsLoopback())

	return endpoint.Scheme == "http" && loopback
}

// environmentBaseURLRefused reports whether ANTHROPIC_BASE_URL, which the SDK
// reads when the client is built, is set to a URL that may not carry the
// key. Checked first, it fails at once instead of after the SDK's retries,
// and a malformed value is refused without the quote the SDK's own parse
// error would give it, userinfo included.
func environmentBaseURLRefused() bool {
	raw := os.Getenv("ANTHROPIC_BASE_URL")
	if raw == "" {
		return false
	}

	endpoint, err := url.Parse(raw)

	return err != nil || !secure(endpoint)
}
