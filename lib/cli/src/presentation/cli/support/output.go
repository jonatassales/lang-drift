package support

import (
	"encoding/json"
	"io"
)

// CheckedWriter records the first write error.
type CheckedWriter struct {
	io.Writer
	err error
}

// NewCheckedWriter wraps an output stream.
func NewCheckedWriter(writer io.Writer) *CheckedWriter {
	return &CheckedWriter{
		Writer: writer,
	}
}

// Write forwards data and records the first error.
func (writer *CheckedWriter) Write(data []byte) (int, error) {
	bytesWritten, writeError := writer.Writer.Write(data)

	if writeError != nil && writer.err == nil {
		writer.err = writeError
	}

	return bytesWritten, writeError
}

// Err returns the first recorded write error.
func (writer *CheckedWriter) Err() error {
	return writer.err
}

// Render writes a successful result as one JSON document: indented for
// people, compact for scripts when --json is set.
func Render(writer io.Writer, result any, compact bool) error {
	encoder := json.NewEncoder(writer)
	if !compact {
		encoder.SetIndent("", "  ")
	}

	return encoder.Encode(result)
}
