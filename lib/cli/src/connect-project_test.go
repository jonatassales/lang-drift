package main

import (
	"context"
	"testing"

	projectdomain "github.com/jonatassales/tao/lib/cli/src/domain/project"
	connectedinfra "github.com/jonatassales/tao/lib/cli/src/infrastructure/connected"
)

// testRuntimeSecret is the deployment secret the composed tests connect with.
const testRuntimeSecret = "composed-test-runtime-secret-0224"

// connectProject records the console a seeded project's planner asks, through
// the store dft agent connect writes with, without the connection check.
func connectProject(t *testing.T, root string, consoleURL string) {
	t.Helper()

	owner := projectdomain.Project{Root: root}

	err := connectedinfra.JudgeConnections{}.Record(context.Background(), owner, consoleURL, testRuntimeSecret)
	if err != nil {
		t.Fatal(err)
	}
}
