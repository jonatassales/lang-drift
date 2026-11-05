package configuration

import (
	"fmt"
	"strings"
	"unicode"
)

// maximumEnvValue bounds one value of a committed env file.
const maximumEnvValue = 4 << 10

// ParseEnvFile reads a committed env file: one NAME=value per line, blank
// lines and lines starting with # ignored. A value is taken as written,
// without quotes or expansion, so the file says exactly what reaches the
// container.
func ParseEnvFile(content string) (map[string]string, error) {
	variables := map[string]string{}

	for number, line := range strings.Split(content, "\n") {
		line = strings.TrimSuffix(line, "\r")
		trimmed := strings.TrimSpace(line)
		blankOrComment := trimmed == "" || strings.HasPrefix(trimmed, "#")

		if blankOrComment {
			continue
		}

		name, value, found := strings.Cut(line, "=")
		if !found {
			return nil, fmt.Errorf("%w: line %d is not NAME=value", ErrInvalidVariable, number+1)
		}

		if err := CheckVariableName(name); err != nil {
			return nil, fmt.Errorf("line %d: %w", number+1, err)
		}

		if _, repeated := variables[name]; repeated {
			return nil, fmt.Errorf("%w: %s is named twice", ErrInvalidVariable, name)
		}

		if len(value) > maximumEnvValue {
			return nil, fmt.Errorf("%w: %s is longer than 4 KiB", ErrInvalidVariable, name)
		}

		if strings.IndexFunc(value, unicode.IsControl) != -1 {
			return nil, fmt.Errorf("%w: %s is not plain text", ErrInvalidVariable, name)
		}

		variables[name] = value
	}

	return variables, nil
}
