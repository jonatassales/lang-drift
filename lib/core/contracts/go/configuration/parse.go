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

	block, err := langdriftBlock(markdown)
	if err != nil {
		return config, err
	}

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
