package main

import (
	"fmt"
	"sort"
)

// goldenDiff compares the golden session with the printed session, aligning
// lines by id and returning every disagreement as a readable message. A line
// listed in pending may differ until its owning spec lands, without hiding a
// kind regression; an unlisted line that differs, or a listed line that
// already matches, is a failure so the pending list only shrinks.
func goldenDiff(golden []goldenLine, printed []printedLine, pending map[string]string) []string {
	var differences []string

	seenPending := map[string]bool{}
	displacedGolden := map[int]bool{}

	printedOffset := func(id string, start int) int {
		for offset, line := range printed[start:] {
			if line.id == id {
				return offset
			}
		}

		return -1
	}

	goldenOffset := func(id string, start int) int {
		for offset, line := range golden[start:] {
			if line.id == id && !displacedGolden[start+offset] {
				return offset
			}
		}

		return -1
	}

	comparePair := func(lineNumber int, expected goldenLine, actual printedLine) {
		spec, isPending := pending[expected.id]
		if isPending {
			seenPending[expected.id] = true
		}

		kindsMatch := expected.kind == actual.kind
		textsMatch := expected.text == actual.text

		if !kindsMatch {
			differences = append(differences, fmt.Sprintf("line %d %q: kind differs: got %s, want %s", lineNumber, expected.id, actual.kind, expected.kind))

			return
		}

		if isPending && textsMatch {
			differences = append(differences, fmt.Sprintf("line %d %q: already matches golden, remove from pending under spec %s", lineNumber, expected.id, spec))

			return
		}

		if !isPending && !textsMatch {
			differences = append(differences, fmt.Sprintf("line %d %q: differs from golden: got %q, want %q", lineNumber, expected.id, actual.text, expected.text))
		}
	}

	recordExtra := func(lineNumber int, actual printedLine, goldenIndex int) {
		_, isPending := pending[actual.id]
		if isPending {
			seenPending[actual.id] = true
		}

		inGolden := goldenOffset(actual.id, goldenIndex)
		if inGolden >= 0 {
			targetIndex := goldenIndex + inGolden
			target := golden[targetIndex]
			if actual.kind != target.kind {
				differences = append(differences, fmt.Sprintf("line %d %q: kind differs: got %s, want %s", lineNumber, actual.id, actual.kind, target.kind))
			}

			if isPending {
				displacedGolden[targetIndex] = true

				return
			}
		}

		if isPending {
			return
		}

		differences = append(differences, fmt.Sprintf("extra line %d %q: got %q", lineNumber, actual.id, actual.text))
	}

	recordMissing := func(lineNumber int, expected goldenLine) {
		if _, isPending := pending[expected.id]; isPending {
			seenPending[expected.id] = true

			return
		}

		differences = append(differences, fmt.Sprintf("missing line %d %q: want %q", lineNumber, expected.id, expected.text))
	}

	goldenIndex := 0
	printedIndex := 0

	for goldenIndex < len(golden) && printedIndex < len(printed) {
		if displacedGolden[goldenIndex] {
			goldenIndex++
			continue
		}

		expected := golden[goldenIndex]
		actual := printed[printedIndex]

		if expected.id == actual.id {
			comparePair(goldenIndex+1, expected, actual)
			goldenIndex++
			printedIndex++
			continue
		}

		inPrinted := printedOffset(expected.id, printedIndex+1)
		inGolden := goldenOffset(actual.id, goldenIndex+1)
		expectedAppearsLaterInPrinted := inPrinted >= 0 && (inGolden < 0 || inPrinted < inGolden)

		switch {
		case expectedAppearsLaterInPrinted:
			recordExtra(printedIndex+1, actual, goldenIndex+1)
			printedIndex++
		case inGolden >= 0:
			recordMissing(goldenIndex+1, expected)
			goldenIndex++
		case pending[expected.id] != "":
			recordMissing(goldenIndex+1, expected)
			goldenIndex++
		case pending[actual.id] != "":
			recordExtra(printedIndex+1, actual, goldenIndex+1)
			printedIndex++
		default:
			recordMissing(goldenIndex+1, expected)
			recordExtra(printedIndex+1, actual, goldenIndex+1)
			goldenIndex++
			printedIndex++
		}
	}

	for ; goldenIndex < len(golden); goldenIndex++ {
		if displacedGolden[goldenIndex] {
			continue
		}

		recordMissing(goldenIndex+1, golden[goldenIndex])
	}

	for ; printedIndex < len(printed); printedIndex++ {
		recordExtra(printedIndex+1, printed[printedIndex], len(golden))
	}

	var unusedPending []string
	for id := range pending {
		if !seenPending[id] {
			unusedPending = append(unusedPending, id)
		}
	}

	sort.Strings(unusedPending)

	for _, id := range unusedPending {
		differences = append(differences, fmt.Sprintf("line %q: already matches golden, remove from pending under spec %s", id, pending[id]))
	}

	return differences
}
