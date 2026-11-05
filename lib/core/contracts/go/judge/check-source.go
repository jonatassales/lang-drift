package judge

// verdicts are the verdicts a source may receive.
var verdicts = []string{"kept", "discarded"}

// checkSource refuses a source without a path or a reason, with a verdict
// outside the set or a confidence outside 0 to 1.
func checkSource(path string, source SourceVerdict) error {
	if err := requireText(path+".path", source.Path); err != nil {
		return err
	}
	if err := requireOneOf(path+".verdict", source.Verdict, verdicts); err != nil {
		return err
	}
	if err := requireText(path+".reason", source.Reason); err != nil {
		return err
	}

	return requireConfidence(path+".confidence", source.Confidence)
}
