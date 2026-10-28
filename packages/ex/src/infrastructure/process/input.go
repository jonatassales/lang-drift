package process

import (
	"context"
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

// inputPollMilliseconds is how often a blocked read rechecks cancellation.
const inputPollMilliseconds = 100

// Input makes the CLI's exclusively owned standard input cancellation-aware,
// so an interrupt ends a blocked prompt. It holds a context because io.Reader
// cannot take one. Polling happens on the caller's goroutine, so no reader
// goroutine outlives a canceled prompt. It neither closes nor changes File.
type Input struct {
	Context context.Context
	File    *os.File
}

// Read waits for input or cancellation. The file must have a single reader.
func (input Input) Read(buffer []byte) (int, error) {
	if len(buffer) == 0 {
		return 0, nil
	}

	descriptors := []unix.PollFd{{Fd: int32(input.File.Fd()), Events: unix.POLLIN}}

	for {
		if err := input.Context.Err(); err != nil {
			return 0, err
		}

		ready, err := unix.Poll(descriptors, inputPollMilliseconds)

		switch {
		case err == unix.EINTR:
			continue
		case err != nil:
			return 0, fmt.Errorf("wait for input: %w", err)
		case ready == 0:
			continue
		case descriptors[0].Revents&unix.POLLNVAL != 0:
			return 0, os.ErrClosed
		}

		return input.File.Read(buffer)
	}
}
