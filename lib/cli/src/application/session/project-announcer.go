package session

// ProjectAnnouncer is a sink that wants to know which project a session
// opened, as a terminal names it in its title.
type ProjectAnnouncer interface {
	Opened(project string)
}
