package judge_test

import (
	"net/http/httptest"

	"github.com/jonatassales/tao/lib/cli/src/application/plan"
)

// consoleSecret is the deployment secret the adapter's tests connect with.
const consoleSecret = "console-test-runtime-secret"

// connectionTo is a recorded connection to the test's console.
func connectionTo(server *httptest.Server) plan.Connection {
	return plan.Connection{
		Recorded: true,
		URL:      server.URL,
		Secret:   consoleSecret,
	}
}
