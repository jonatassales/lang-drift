package configuration

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

// versioned reads only the fence's version, before its shape is trusted.
type versioned struct {
	Version int `json:"version"`
}

// Parse requires one explicit langdrift JSON fence, rejects ambiguous or unknown
// fields, and bounds the document. Markdown never grants execution authority.
func Parse(markdown string) (Configuration, error) {
	var config Configuration
	if len(markdown) > 1<<20 {
		return config, invalid("Markdown exceeds 1 MiB")
	}

	const fence = "```langdrift\n"
	if strings.Count(markdown, fence) != 1 {
		return config, invalid("Markdown requires exactly one langdrift JSON fence")
	}

	after := strings.SplitN(markdown, fence, 2)[1]
	end := strings.Index(after, "\n```")
	if end < 0 {
		return config, invalid("configuration fence is not closed")
	}

	block := after[:end]

	var header versioned
	if err := json.Unmarshal([]byte(block), &header); err != nil {
		return config, fmt.Errorf("configuration block: %w", err)
	}

	if header.Version != Version {
		return config, invalid(fmt.Sprintf("version %d is not version %d; %s", header.Version, Version, replan))
	}

	decoder := json.NewDecoder(strings.NewReader(block))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&config); err != nil {
		return config, invalid(fmt.Sprintf("%v; %s", err, replan))
	}

	if err := decoder.Decode(new(any)); !errors.Is(err, io.EOF) {
		return config, invalid("configuration block must contain one JSON object")
	}

	return config, config.Validate()
}
