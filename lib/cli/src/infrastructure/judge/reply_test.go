package judge_test

import (
	"net/http"
	"testing"
)

// reply answers with a status and a body, failing the test when the body
// cannot be written.
func reply(t *testing.T, writer http.ResponseWriter, status int, body string) {
	t.Helper()

	writer.WriteHeader(status)

	if _, err := writer.Write([]byte(body)); err != nil {
		t.Error(err)
	}
}
