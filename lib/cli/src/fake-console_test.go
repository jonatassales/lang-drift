package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	contract "github.com/jonatassales/tao/lib/core/contracts/go/judge"
)

// fakeConsole plays the console's judge route in-process: the connection
// check answers 204 to the secret and 401 without it, and each question is
// answered from the scripted texts, or with the status and body it is told,
// after the delay it is told. It remembers every question and bearer.
type fakeConsole struct {
	server  *httptest.Server
	test    *testing.T
	answers map[string]string
	status  int
	body    string
	refusal string
	delay   time.Duration
	mutex   sync.Mutex
	asked   []contract.Request
	bearers []string
}

func (console *fakeConsole) serve(writer http.ResponseWriter, request *http.Request) {
	authorized := request.Header.Get("Authorization") == "Bearer "+testRuntimeSecret

	if request.Method == http.MethodGet {
		console.check(writer, authorized)

		return
	}

	question := console.record(request)

	if console.delay > 0 {
		select {
		case <-time.After(console.delay):
		case <-request.Context().Done():
			return
		}
	}

	if console.status != 0 {
		writer.WriteHeader(console.status)
		console.write(writer, []byte(console.body))

		return
	}

	if !authorized {
		writer.WriteHeader(http.StatusUnauthorized)

		return
	}

	console.answer(writer, question)
}

func (console *fakeConsole) check(writer http.ResponseWriter, authorized bool) {
	if !authorized {
		writer.WriteHeader(http.StatusUnauthorized)

		return
	}

	writer.WriteHeader(http.StatusNoContent)
}

// record reads one question strictly as the contract's request.
func (console *fakeConsole) record(request *http.Request) contract.Request {
	var question contract.Request

	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&question); err != nil {
		console.test.Errorf("the judge route received a body outside the contract: %v", err)
	}

	console.mutex.Lock()
	defer console.mutex.Unlock()

	console.asked = append(console.asked, question)
	console.bearers = append(console.bearers, request.Header.Get("Authorization"))

	return question
}

func (console *fakeConsole) answer(writer http.ResponseWriter, question contract.Request) {
	response := contract.Response{
		Text: console.answers[string(question.Kind)],
		Usage: contract.Usage{
			InputTokens:  15000,
			OutputTokens: 1500,
		},
		Cost: contract.Cost{
			Cents:        5,
			Provenance:   "console",
			Illustrative: true,
		},
		JudgeVersion: "console-judge-1",
	}

	if console.refusal != "" {
		response.Text = ""
		response.Refused = true
		response.Reason = console.refusal
	}

	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(response); err != nil {
		console.test.Error(err)
	}

	writer.Header().Set("Content-Type", "application/json")
	console.write(writer, body.Bytes())
}

// write sends a body, failing the test when it cannot.
func (console *fakeConsole) write(writer http.ResponseWriter, body []byte) {
	if _, err := writer.Write(body); err != nil {
		console.test.Error(err)
	}
}

// questions is what the console was asked, and with which bearers.
func (console *fakeConsole) questions() ([]contract.Request, []string) {
	console.mutex.Lock()
	defer console.mutex.Unlock()

	return append([]contract.Request{}, console.asked...), append([]string{}, console.bearers...)
}
