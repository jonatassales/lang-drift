package main

import (
	"flag"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
)

// updateExpected rewrites the expected tree from the session the test
// performs: go test -run TestSessionWritesTheExpectedTree -update .
var updateExpected = flag.Bool("update", false, "rewrite testdata/langdrift/expected from the session")

const expectedFixture = "testdata/langdrift/expected"

// privateState is the part of the folder the fixture leaves out: the laps'
// state, the index and the operator records are the CLI's own.
const privateState = "agent/state"

// The session, performed end to end with scripted sources, writes exactly
// the tree under the expected fixture: the records are the contract, byte
// for byte, except the private state.
func TestSessionWritesTheExpectedTree(t *testing.T) {
	t.Parallel()

	p, deps := seedProject(t)
	deps.Judge = &scriptedJudge{answers: sampleAnswers(t)}
	deps.Driver, deps.Lifecycle, deps.Delivery, deps.Poll = websiteDriver(), websiteLifecycle(), websiteDelivery(), time.Millisecond

	if status, stdout, stderr := p.dftWith(deps, "agent", "plan", consoleTicket, "--component", consoleComponent); status != 0 {
		t.Fatalf("dft agent plan = %d: %s\n%s", status, stderr, stdout)
	}
	if status, stdout, stderr := p.dftWith(deps, "agent", "run", "console-microruntime", "--model", "claude-opus-5-5", "--effort", "high"); status != 3 {
		t.Fatalf("the first lap gave %d: %s\n%s", status, stderr, stdout)
	}
	if status, stdout, stderr := p.dftInput(deps, "y\n", "agent", "run", "console-microruntime", "--confirm"); status != 0 {
		t.Fatalf("the second lap gave %d: %s\n%s", status, stderr, stdout)
	}

	actual := treeOf(t, filepath.Join(p.root, ".langdrift"))

	if *updateExpected {
		rewriteExpected(t, actual)
	}

	expected := treeOf(t, expectedFixture)

	for _, name := range sortedKeys(expected) {
		content, found := actual[name]
		switch {
		case !found:
			t.Errorf("%s is expected but the session did not write it", name)
		case content != expected[name]:
			t.Errorf("%s differs from the expected fixture:\n%s", name, lineDifference(expected[name], content))
		}
	}
	for _, name := range sortedKeys(actual) {
		if _, found := expected[name]; !found {
			t.Errorf("the session wrote %s, which the expected fixture does not hold", name)
		}
	}
}

// treeOf reads every file under a folder, keyed by its slash path, leaving
// the private state out.
func treeOf(t *testing.T, root string) map[string]string {
	t.Helper()

	tree := map[string]string{}
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		name := filepath.ToSlash(relative)
		if entry.IsDir() {
			if name == privateState {
				return filepath.SkipDir
			}

			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		tree[name] = string(content)

		return nil
	})
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}

	return tree
}

// rewriteExpected replaces the expected fixture with the session's tree.
func rewriteExpected(t *testing.T, tree map[string]string) {
	t.Helper()

	if err := os.RemoveAll(expectedFixture); err != nil {
		t.Fatal(err)
	}
	for name, content := range tree {
		path := filepath.Join(expectedFixture, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	t.Logf("rewrote %d files under %s", len(tree), expectedFixture)
}

func sortedKeys(tree map[string]string) []string {
	keys := make([]string, 0, len(tree))
	for key := range tree {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	return keys
}

// firstDifference shows the first line that differs, with its number.
func lineDifference(expected, actual string) string {
	expectedLines, actualLines := strings.Split(expected, "\n"), strings.Split(actual, "\n")
	for index := range max(len(expectedLines), len(actualLines)) {
		var want, got string
		if index < len(expectedLines) {
			want = expectedLines[index]
		}
		if index < len(actualLines) {
			got = actualLines[index]
		}
		if want != got {
			return "line " + strconv.Itoa(index+1) + "\nexpected: " + want + "\nactual:   " + got
		}
	}

	return "no line differs"
}
