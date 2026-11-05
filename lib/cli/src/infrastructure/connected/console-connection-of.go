package connected

import (
	"errors"
	"strings"

	"github.com/jonatassales/tao/lib/cli/src/application/plan"
)

// errIncompleteConnection refuses a console connection file that lacks the URL
// or the secret, which only an edit by hand or a crash leaves behind.
var errIncompleteConnection = errors.New("the console connection is incomplete; run dft agent connect <console-url> again")

// consoleConnectionOf reads the private console connection file: one
// CONSOLE_URL and one RUNTIME_SECRET line.
func consoleConnectionOf(content string) (plan.Connection, error) {
	values := map[string]string{}

	for _, line := range strings.Split(content, "\n") {
		name, value, found := strings.Cut(line, "=")
		if found {
			values[name] = value
		}
	}

	connection := plan.Connection{
		Recorded: true,
		URL:      values["CONSOLE_URL"],
		Secret:   values["RUNTIME_SECRET"],
	}

	if connection.URL == "" || connection.Secret == "" {
		return plan.Connection{}, errIncompleteConnection
	}

	return connection, nil
}
