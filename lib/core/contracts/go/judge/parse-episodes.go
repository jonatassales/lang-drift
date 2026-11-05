package judge

import "encoding/json"

// The fields of an episodes answer and of each episode.
var (
	episodesFields = []string{"episodes"}
	episodeFields  = []string{"run", "reason"}
)

// rawEpisodes is an episodes answer before its episodes are read.
type rawEpisodes struct {
	Episodes json.RawMessage `json:"episodes"`
}

// parseEpisodes reads an episodes answer: each episode names its run and the
// reason it is or is not comparable.
func parseEpisodes(text []byte) (*EpisodesAnswer, error) {
	var raw rawEpisodes
	if err := decodeObject(text, "", episodesFields, &raw); err != nil {
		return nil, err
	}

	episodes, err := decodeItems(raw.Episodes, "episodes", episodeFields, checkEpisode)
	if err != nil {
		return nil, err
	}

	return &EpisodesAnswer{Episodes: episodes}, nil
}
