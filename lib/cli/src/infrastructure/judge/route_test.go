package judge_test

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

// route plays the console's judge route with one handler, counting the calls
// it receives.
func route(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *atomic.Int32) {
	t.Helper()

	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		calls.Add(1)
		handler(writer, request)
	}))
	t.Cleanup(server.Close)

	return server, &calls
}
