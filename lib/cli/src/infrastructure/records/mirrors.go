package records

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"strconv"

	"github.com/jonatassales/tao/lib/cli/src/domain/project"
	"github.com/jonatassales/tao/lib/cli/src/infrastructure/files"
)

// documentMode is the mode of a record or mirror: the consumer reads it.
const documentMode = 0o644

var numberedPrefix = regexp.MustCompile(`^(\d{4})-`)

// Mirrors writes the mirrors of the inputs into the .langdrift folder.
type Mirrors struct{}

// WriteMirror writes one mirror, relative to the folder, in the file layout.
func (Mirrors) WriteMirror(owner project.Project, relative string, value any) error {
	target, err := files.Within(owner.Root, ".langdrift/"+relative)
	if err != nil {
		return err
	}

	data, err := Format(value, FileProfile)
	if err != nil {
		return err
	}

	return files.WriteDocument(target, data, documentMode)
}

// NextPriceNumber is the next free four-digit number under finops/prices.
func (Mirrors) NextPriceNumber(owner project.Project) (string, error) {
	directory, err := files.Within(owner.Root, ".langdrift/finops/prices")
	if err != nil {
		return "", err
	}

	entries, err := os.ReadDir(directory)
	if errors.Is(err, fs.ErrNotExist) {
		return "0001", nil
	}
	if err != nil {
		return "", err
	}

	highest := 0
	for _, entry := range entries {
		match := numberedPrefix.FindStringSubmatch(entry.Name())
		if match == nil {
			continue
		}

		number, _ := strconv.Atoi(match[1])
		highest = max(highest, number)
	}

	if highest >= 9999 {
		return "", errors.New("no price list number is left")
	}

	return fmt.Sprintf("%04d", highest+1), nil
}
