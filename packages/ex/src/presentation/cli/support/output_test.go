package support_test

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/exphere/ex/src/presentation/cli/support"
)

type failingWriter struct{ err error }

func (writer failingWriter) Write([]byte) (int, error) { return 0, writer.err }

func TestCheckedWriter(t *testing.T) {
	t.Parallel()

	var sink strings.Builder
	writer := support.NewCheckedWriter(&sink)

	if _, err := io.WriteString(writer, "ok"); err != nil || writer.Err() != nil || sink.String() != "ok" {
		t.Fatalf("a successful write = %v, recorded %v", err, writer.Err())
	}

	first, second := errors.New("broken pipe"), errors.New("closed")
	failing := support.NewCheckedWriter(failingWriter{err: first})
	_, _ = failing.Write([]byte("a"))
	failing.Writer = failingWriter{err: second}
	_, _ = failing.Write([]byte("b"))

	if !errors.Is(failing.Err(), first) {
		t.Fatalf("Err() = %v, want the first write error", failing.Err())
	}
}

func TestRender(t *testing.T) {
	t.Parallel()

	result := map[string]any{"name": "demo", "units": []string{"web"}}

	var indented, compact strings.Builder
	if err := support.Render(&indented, result, false); err != nil {
		t.Fatal(err)
	}

	if err := support.Render(&compact, result, true); err != nil {
		t.Fatal(err)
	}

	if indented.String() != "{\n  \"name\": \"demo\",\n  \"units\": [\n    \"web\"\n  ]\n}\n" {
		t.Errorf("indented = %q", indented.String())
	}

	if compact.String() != "{\"name\":\"demo\",\"units\":[\"web\"]}\n" {
		t.Errorf("compact = %q", compact.String())
	}
}
