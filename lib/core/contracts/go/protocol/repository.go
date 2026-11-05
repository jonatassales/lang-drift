package protocol

// Repository names the source a runtime builds from: a local folder or a
// GitHub repository at a revision.
type Repository struct {
	Kind     string `json:"kind,omitempty"`
	Path     string `json:"path,omitempty"`
	FullName string `json:"fullName"`
	URL      string `json:"url"`
	Ref      string `json:"ref"`
}
