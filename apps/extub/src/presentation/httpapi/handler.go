// Package httpapi exposes the local seed contracts without native or provider access.
package httpapi

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/jonatassales/tao/apps/extub/src/application"
	"github.com/jonatassales/tao/apps/extub/src/domain"
)

type Handler struct {
	Service  application.Service
	Signaler application.Signaler
	// Token switches extub to its hosted preview: empty keeps the local,
	// loopback-only contract; set, every request but the health check must
	// carry it, and nothing may change the seed.
	Token string
}

// hosted is true when extub runs as the token-gated preview.
func (handler Handler) hosted() bool {
	return handler.Token != ""
}

// carriesToken reads a Bearer token (RFC 6750), its scheme in any case as
// RFC 9110 has it, and compares SHA-256 digests in constant time, so the
// comparison reveals neither the token nor its length.
func (handler Handler) carriesToken(request *http.Request) bool {
	scheme, token, found := strings.Cut(request.Header.Get("Authorization"), " ")
	if !found || !strings.EqualFold(scheme, "Bearer") {
		return false
	}
	presented := sha256.Sum256([]byte(strings.TrimLeft(token, " ")))
	expected := sha256.Sum256([]byte(handler.Token))
	return subtle.ConstantTimeCompare(presented[:], expected[:]) == 1
}

// admit applies the hosted preview's rules: reading the health check needs
// nothing, every other request needs the token, and only reads pass. A
// WebRTC offer passes to its route, which answers that the preview has no
// stream.
func (handler Handler) admit(writer http.ResponseWriter, request *http.Request) bool {
	path := strings.Trim(request.URL.Path, "/")
	read := request.Method == http.MethodGet
	if path == "health" && read {
		return true
	}
	if !handler.carriesToken(request) {
		writer.Header().Set("WWW-Authenticate", `Bearer realm="extub"`)
		respondError(writer, http.StatusUnauthorized, "token_required")
		return false
	}
	if !read && path != "webrtc/offer" {
		respondError(writer, http.StatusForbidden, "read_only_preview")
		return false
	}
	return true
}

func (handler Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	setAPIHeaders(writer.Header())
	if handler.hosted() {
		if handler.admit(writer, request) {
			handler.respond(writer, request)
		}
		return
	}
	host := request.Host
	if name, _, err := net.SplitHostPort(host); err == nil {
		host = name
	}
	origin := request.Header.Get("Origin")
	if (host != "localhost" && host != "127.0.0.1") || (origin != "" && origin != "http://"+request.Host) || request.Header.Get("Sec-Fetch-Site") == "cross-site" {
		respondError(writer, http.StatusForbidden, "local_only")
		return
	}
	if request.URL.Path == "/probe" && request.Method == http.MethodGet {
		writer.WriteHeader(http.StatusNoContent)
		return
	}
	handler.respond(writer, request)
}

// setAPIHeaders keeps every answer data for the console alone: JSON that is
// never stored, sniffed, rendered, framed, sent as a referrer or read by
// another origin.
func setAPIHeaders(header http.Header) {
	header.Set("Content-Type", "application/json")
	header.Set("Cache-Control", "no-store")
	header.Set("X-Content-Type-Options", "nosniff")
	header.Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'")
	header.Set("X-Frame-Options", "DENY")
	header.Set("Referrer-Policy", "no-referrer")
	header.Set("Cross-Origin-Resource-Policy", "same-origin")
}

// bodyLimit bounds a request body at 64 KiB.
const bodyLimit = 64 << 10

// respond routes an admitted request and writes its result or its error.
func (handler Handler) respond(writer http.ResponseWriter, request *http.Request) {
	request.Body = http.MaxBytesReader(writer, request.Body, bodyLimit)
	result, err := handler.route(request)
	if err != nil {
		status := http.StatusInternalServerError
		code := "operation_failed"
		switch {
		case errors.Is(err, domain.ErrNotFound):
			status = http.StatusNotFound
			code = "not_found"
		case errors.Is(err, domain.ErrInvalid):
			status = http.StatusBadRequest
			code = "invalid_input"
		case errors.Is(err, errTooLarge):
			status = http.StatusRequestEntityTooLarge
			code = "request_too_large"
		case errors.Is(err, domain.ErrConflict):
			status = http.StatusConflict
			code = "executions_not_stopped"
		case errors.Is(err, domain.ErrLimit):
			status = http.StatusTooManyRequests
			code = "seed_limit"
		case errors.Is(err, domain.ErrLoadTestBusy):
			status = http.StatusTooManyRequests
			code = "load_test_busy"
		case errors.Is(err, errMethod):
			status = http.StatusMethodNotAllowed
			code = "method_not_allowed"
		case errors.Is(err, errStreamUnavailable):
			status = http.StatusServiceUnavailable
			code = "stream_unavailable"
		case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
			status = http.StatusRequestTimeout
			code = "negotiation_timeout"
		}
		respondError(writer, status, code)
		return
	}
	_ = json.NewEncoder(writer).Encode(result)
}

var errMethod = errors.New("method_not_allowed")

// errStreamUnavailable answers a WebRTC offer the hosted preview cannot take:
// Vercel's functions receive no UDP, so the console keeps polling instead.
var errStreamUnavailable = errors.New("stream_unavailable")

// errTooLarge answers a body over bodyLimit: refused before routing when its
// length is declared, and when a route reads past the limit otherwise.
var errTooLarge = errors.New("request_too_large")

func respondError(writer http.ResponseWriter, status int, code string) {
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(map[string]any{"source": "stub", "error": map[string]string{"code": code}})
}

// decode reads exactly one JSON value from the body.
func decode(request *http.Request, target any) error {
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(target); err != nil {
		return decodeFailure(err)
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		return decodeFailure(err)
	}
	return nil
}

// decodeFailure tells a body that read past the limit from one that is not
// exactly one JSON value.
func decodeFailure(err error) error {
	var tooLarge *http.MaxBytesError
	if errors.As(err, &tooLarge) {
		return errTooLarge
	}
	return domain.ErrInvalid
}

func envelope(key string, value any) map[string]any {
	return map[string]any{"source": "stub", key: value}
}
func (handler Handler) route(request *http.Request) (any, error) {
	if request.ContentLength > bodyLimit {
		return nil, errTooLarge
	}
	path := strings.Trim(request.URL.Path, "/")
	method := request.Method
	if method == http.MethodGet && path == "health" {
		return map[string]any{"source": "stub", "service": "extub", "status": "ready"}, nil
	}
	if path == "webrtc/offer" {
		if method != http.MethodPost {
			return nil, errMethod
		}
		if handler.hosted() {
			return nil, errStreamUnavailable
		}
		var input application.Description
		if err := decode(request, &input); err != nil {
			return nil, err
		}
		return handler.Signaler.Offer(request.Context(), input)
	}
	if !strings.HasPrefix(path, "console/") {
		return nil, domain.ErrNotFound
	}
	path = strings.TrimPrefix(path, "console/")
	if path == "meta" && method == http.MethodGet {
		return map[string]any{"source": "stub", "service": "extub", "webrtc": true}, nil
	}
	if strings.HasPrefix(path, "auth/") {
		if (path == "auth/session" && method != http.MethodGet) || (path != "auth/session" && method != http.MethodPost) {
			return nil, errMethod
		}
		return envelope("user", map[string]string{"id": "seed-user", "email": "preview@example.test"}), nil
	}
	if path == "local/folders" && method == http.MethodGet {
		folder := request.URL.Query().Get("path")
		if folder == "" {
			folder = "/seeds"
		}
		if err := domain.ValidateFolder(folder); err != nil {
			return nil, err
		}
		return map[string]any{"source": "stub", "path": folder, "parent": nil, "folders": []map[string]string{{"name": "microruntimes", "path": "/seeds/microruntimes"}}}, nil
	}
	if path == "home" && method == http.MethodGet {
		return handler.home(), nil
	}
	parts := strings.Split(path, "/")
	switch parts[0] {
	case "products":
		return handler.products(request, parts[1:])
	case "runtimes":
		return handler.runtimes(request, parts[1:])
	case "tickets":
		return handler.tickets(request, parts[1:])
	}
	return nil, domain.ErrNotFound
}
func (handler Handler) home() any {
	catalog := handler.Service.Snapshot()
	running := []any{}
	unknown := 0
	for _, widget := range catalog.Widgets(handler.Service.Now(), 0) {
		active := false
		for _, execution := range widget.Executions {
			if execution.State == "running" {
				active = true
			}
			if execution.State == "unknown" {
				unknown++
			}
		}
		if active {
			product, _ := catalog.Product(widget.Runtime.ProductID)
			running = append(running, map[string]any{"runtime": widget.Runtime, "productName": product.Name, "executions": widget.Executions})
		}
	}
	return map[string]any{"source": "stub", "products": catalog.Products, "running": running, "unknownCount": unknown}
}
