package judge

// checkEpisode refuses an episode without its run or its reason.
func checkEpisode(path string, episode EpisodeReason) error {
	if err := requireText(path+".run", episode.Run); err != nil {
		return err
	}

	return requireText(path+".reason", episode.Reason)
}
