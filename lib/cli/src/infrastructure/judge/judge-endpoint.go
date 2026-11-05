package judge

import (
	"errors"
	"net/url"
	"strings"
)

// judgePath is the console's judge route.
const judgePath = "/api/judge"

// errMalformedURL refuses a URL that does not parse, without quoting it: the
// parse error would repeat its userinfo.
var errMalformedURL = errors.New("the console URL is malformed")

// errURLParts refuses a URL that carries credentials, a query or a fragment.
var errURLParts = errors.New("the console URL must name a host, without credentials, query or fragment")

// judgeEndpoint is the judge route of a console URL: HTTPS, or HTTP on
// loopback, with no credentials, query or fragment of its own.
func judgeEndpoint(raw string) (string, error) {
	endpoint, err := url.Parse(raw)
	if err != nil {
		return "", errMalformedURL
	}

	extraParts := endpoint.User != nil || endpoint.RawQuery != "" || endpoint.Fragment != ""
	if extraParts || endpoint.Hostname() == "" {
		return "", errURLParts
	}

	if !secure(endpoint) {
		return "", errInsecureURL
	}

	endpoint.Path = strings.TrimSuffix(strings.TrimRight(endpoint.Path, "/"), judgePath) + judgePath

	return endpoint.String(), nil
}
