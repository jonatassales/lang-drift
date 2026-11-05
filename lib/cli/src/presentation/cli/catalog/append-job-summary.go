package catalog

import (
	"errors"
	"fmt"
	"os"
)

// appendJobSummary appends text to the job summary file GitHub Actions names,
// after the summaries of the job's earlier steps.
func appendJobSummary(path, text string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o644)
	if err != nil {
		return fmt.Errorf("open the job summary: %w", err)
	}

	if _, err := fmt.Fprint(file, text); err != nil {
		return fmt.Errorf("write the job summary: %w", errors.Join(err, file.Close()))
	}

	if err := file.Close(); err != nil {
		return fmt.Errorf("write the job summary: %w", err)
	}

	return nil
}
