package manifest

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// commandOf reads the CMD of a HEALTHCHECK in its exec or its shell form.
func commandOf(command string) ([]string, error) {
	if command == "" {
		return nil, errors.New("HEALTHCHECK CMD is empty")
	}

	if !strings.HasPrefix(command, "[") {
		return []string{"/bin/sh", "-c", command}, nil
	}

	var arguments []string
	if err := json.Unmarshal([]byte(command), &arguments); err != nil || len(arguments) == 0 {
		return nil, fmt.Errorf("HEALTHCHECK CMD %s is not a JSON array of strings", command)
	}

	return arguments, nil
}
