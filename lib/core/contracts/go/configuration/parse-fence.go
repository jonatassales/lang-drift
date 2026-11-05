package configuration

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

// ParseFence reads the one langdrift fence of a runtime.md as version 3,
// refusing any other version with "re-plan the runtime" and any field the
// version does not define. It reads the shape only: Validate says which
// rules the fence breaks.
func ParseFence(markdown string) (Fence, error) {
	block, err := langdriftBlock(markdown)
	if err != nil {
		return Fence{}, err
	}

	var header versioned
	if err := json.Unmarshal([]byte(block), &header); err != nil {
		return Fence{}, fmt.Errorf("configuration block: %w", err)
	}

	if header.Version != FenceVersion {
		return Fence{}, invalid(fmt.Sprintf("version %d is not version %d; %s", header.Version, FenceVersion, replan))
	}

	decoder := json.NewDecoder(strings.NewReader(block))
	decoder.DisallowUnknownFields()

	var fence Fence
	if err := decoder.Decode(&fence); err != nil {
		return Fence{}, invalid(fmt.Sprintf("%v; %s", err, replan))
	}

	if err := decoder.Decode(new(any)); !errors.Is(err, io.EOF) {
		return Fence{}, invalid("configuration block must contain one JSON object")
	}

	return fence, nil
}
