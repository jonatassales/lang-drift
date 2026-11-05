package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// newFakeConsole starts a console whose judge route answers questions from
// the scripted texts, and stops it when the test ends.
func newFakeConsole(t *testing.T, answers map[string]string) *fakeConsole {
	t.Helper()

	console := &fakeConsole{
		test:    t,
		answers: answers,
	}
	console.server = httptest.NewServer(http.HandlerFunc(console.serve))
	t.Cleanup(console.server.Close)

	return console
}
