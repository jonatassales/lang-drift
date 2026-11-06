package session

// Announce tells a sink which project the session opened, by its folder's
// name, when the sink wants to know; other sinks are left alone.
func Announce(sink Sink, project string) {
	announcer, announces := sink.(ProjectAnnouncer)
	if !announces {
		return
	}

	announcer.Opened(project)
}
