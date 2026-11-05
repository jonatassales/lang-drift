package judgment

import (
	"regexp"
	"slices"
	"strconv"
	"strings"
)

// episodeID is the form of a project episode's id, as ep-0412.
var episodeID = regexp.MustCompile(`^ep-[0-9]{4}$`)

// checkEpisodes accepts the judge's verdicts only as the EPISODES rows can
// show them: every episode once, by its ep-NNNN id, kept because it shares
// the stack and scale, or dropped for one of the closed reasons. A refusal
// names the entry's field and the rule it breaks.
func checkEpisodes(answer EpisodesAnswer) error {
	seen := map[string]bool{}

	for index, episode := range answer.Episodes {
		field := "episodes[" + strconv.Itoa(index) + "]."
		keptFits := episode.Verdict == Kept && episode.Reason == KeptReason
		droppedFits := episode.Verdict == Dropped && slices.Contains(DroppedReasons, episode.Reason)
		knownVerdict := episode.Verdict == Kept || episode.Verdict == Dropped

		switch {
		case !episodeID.MatchString(episode.ID):
			return FieldError{
				Field: field + "id",
				Rule:  "must be an episode id, ep-NNNN",
			}
		case seen[episode.ID]:
			return FieldError{
				Field: field + "id",
				Rule:  "must judge each episode once",
			}
		case !knownVerdict:
			return FieldError{
				Field: field + "verdict",
				Rule:  "must be kept or dropped",
			}
		case !keptFits && !droppedFits:
			return FieldError{
				Field: field + "reason",
				Rule:  "must be the verdict's reason: " + KeptReason + " when kept; " + strings.Join(DroppedReasons, ", ") + " when dropped",
			}
		}

		seen[episode.ID] = true
	}

	return nil
}
