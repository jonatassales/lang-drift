package process_test

import (
	"context"
	"errors"
	"io"
	"os"
	"testing"
	"time"

	"github.com/exphere/ex/src/infrastructure/process"
)

func pipe(t *testing.T) (*os.File, *os.File) {
	t.Helper()

	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		_ = reader.Close()
		_ = writer.Close()
	})

	return reader, writer
}

func TestInput_Read(t *testing.T) {
	t.Parallel()

	t.Run("returns available input", func(t *testing.T) {
		t.Parallel()

		reader, writer := pipe(t)
		if _, err := writer.WriteString("answer\n"); err != nil {
			t.Fatal(err)
		}

		buffer := make([]byte, 64)
		count, err := process.Input{Context: context.Background(), File: reader}.Read(buffer)
		if err != nil || string(buffer[:count]) != "answer\n" {
			t.Fatalf("Read() = %q, %v", buffer[:count], err)
		}
	})

	t.Run("reports the end of input", func(t *testing.T) {
		t.Parallel()

		reader, writer := pipe(t)
		_ = writer.Close()

		if count, err := (process.Input{Context: context.Background(), File: reader}).Read(make([]byte, 8)); count != 0 || !errors.Is(err, io.EOF) {
			t.Fatalf("Read() = %d, %v; want io.EOF", count, err)
		}
	})

	t.Run("an empty buffer reads nothing", func(t *testing.T) {
		t.Parallel()

		reader, _ := pipe(t)
		if count, err := (process.Input{Context: context.Background(), File: reader}).Read(nil); count != 0 || err != nil {
			t.Fatalf("Read(nil) = %d, %v", count, err)
		}
	})

	t.Run("cancellation ends a blocked read", func(t *testing.T) {
		t.Parallel()

		reader, _ := pipe(t)
		ctx, cancel := context.WithCancel(context.Background())
		time.AfterFunc(50*time.Millisecond, cancel)

		started := time.Now()
		_, err := process.Input{Context: ctx, File: reader}.Read(make([]byte, 8))

		if !errors.Is(err, context.Canceled) || time.Since(started) > 2*time.Second {
			t.Fatalf("Read() = %v after %v; want prompt cancellation", err, time.Since(started))
		}
	})
}
