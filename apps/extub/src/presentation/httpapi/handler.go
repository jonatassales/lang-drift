// Package httpapi exposes the local seed contracts without native or provider access.
package httpapi

import (
	"context"
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

// carriesToken compares the bearer token in constant time.
func (handler Handler) carriesToken(request *http.Request) bool {
	presented, found := strings.CutPrefix(request.Header.Get("Authorization"), "Bearer ")
	return found && subtle.ConstantTimeCompare([]byte(presented), []byte(handler.Token)) == 1
}

// readOnly is true for the methods the hosted preview serves.
func readOnly(method string) bool {
	return method == http.MethodGet || method == http.MethodHead
}

// admit applies the hosted preview's rules: the health check is open, every
// other request needs the token, and only reads pass.
func (handler Handler) admit(writer http.ResponseWriter, request *http.Request) bool {
	if strings.Trim(request.URL.Path, "/") == "health" {
		return true
	}
	if !handler.carriesToken(request) {
		respondError(writer, http.StatusUnauthorized, "token_required")
		return false
	}
	if !readOnly(request.Method) && request.URL.Path != "/webrtc/offer" {
		respondError(writer, http.StatusForbidden, "read_only_preview")
		return false
	}
	return true
}

func (handler Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Cache-Control", "no-store")
	writer.Header().Set("X-Content-Type-Options", "nosniff")
	if handler.hosted() {
		request.Body = http.MaxBytesReader(writer, request.Body, 65536)
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
	request.Body = http.MaxBytesReader(writer, request.Body, 65536)
	if request.URL.Path == "/probe" && request.Method == http.MethodGet {
		writer.WriteHeader(http.StatusNoContent)
		return
	}
	handler.respond(writer, request)
}

// respond routes an admitted request and writes its result or its error.
func (handler Handler) respond(writer http.ResponseWriter, request *http.Request) {
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

func respondError(writer http.ResponseWriter, status int, code string) {
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(map[string]any{"source": "stub", "error": map[string]string{"code": code}})
}
func decode(request *http.Request, target any) error {
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(target); err != nil {
		return domain.ErrInvalid
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		return domain.ErrInvalid
	}
	return nil
}
func envelope(key string, value any) map[string]any {
	return map[string]any{"source": "stub", key: value}
}
func (handler Handler) route(request *http.Request) (any, error) {
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
