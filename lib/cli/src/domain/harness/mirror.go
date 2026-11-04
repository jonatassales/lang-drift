package harness

// Mirror opens every mirrored file: which record it mirrors, who set the
// record, through what, and when. No gate reads a mirror; an edit to one is
// shown as an unauthorized change.
type Mirror struct {
	Of       string    `json:"of"`
	Record   string    `json:"record"`
	SetBy    string    `json:"setBy"`
	Via      string    `json:"via"`
	SetAt    Timestamp `json:"setAt"`
	ReadOnly bool      `json:"readOnly"`
	Note     string    `json:"note"`
}

// Resources is a CPU and memory pair, in millicores and MiB.
type Resources struct {
	CPUMillis int `json:"cpuMillis"`
	MemoryMiB int `json:"memoryMiB"`
}
