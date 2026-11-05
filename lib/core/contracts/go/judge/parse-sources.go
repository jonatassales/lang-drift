package judge

import "encoding/json"

// The fields of a sources answer and of each of its items.
var (
	sourcesFields   = []string{"sources", "facts", "criteria", "figures"}
	sourceFields    = []string{"path", "verdict", "reason", "confidence", "role", "summary"}
	factFields      = []string{"topic", "text", "source"}
	criterionFields = []string{"id", "applies"}
	figureFields    = []string{"topic", "fact", "source", "cpuMillis", "memoryMiB"}
)

// rawSources is a sources answer before its items are read.
type rawSources struct {
	Sources  json.RawMessage `json:"sources"`
	Facts    json.RawMessage `json:"facts"`
	Criteria json.RawMessage `json:"criteria"`
	Figures  json.RawMessage `json:"figures"`
}

// parseSources reads a sources answer, each of its four lists item by item.
func parseSources(text []byte) (*SourcesAnswer, error) {
	var raw rawSources
	if err := decodeObject(text, "", sourcesFields, &raw); err != nil {
		return nil, err
	}

	sources, err := decodeItems(raw.Sources, "sources", sourceFields, checkSource)
	if err != nil {
		return nil, err
	}

	facts, err := decodeItems(raw.Facts, "facts", factFields, checkFact)
	if err != nil {
		return nil, err
	}

	criteria, err := decodeItems(raw.Criteria, "criteria", criterionFields, checkCriterion)
	if err != nil {
		return nil, err
	}

	figures, err := decodeItems(raw.Figures, "figures", figureFields, checkFigure)
	if err != nil {
		return nil, err
	}

	return &SourcesAnswer{
		Sources:  sources,
		Facts:    facts,
		Criteria: criteria,
		Figures:  figures,
	}, nil
}
