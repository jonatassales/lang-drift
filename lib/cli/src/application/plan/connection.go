package plan

// Connection is the console the project is connected to, as dft agent
// connect recorded it: the console's URL and the deployment secret its judge
// route accepts. A project that never connected has none recorded.
type Connection struct {
	Recorded bool
	URL      string
	Secret   string
}

// String names the console and never the secret, so a connection printed by
// mistake leaks nothing.
func (connection Connection) String() string {
	if !connection.Recorded {
		return "no console connection"
	}

	return "console " + connection.URL
}
