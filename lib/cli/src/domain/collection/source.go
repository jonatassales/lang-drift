package collection

// Source is one file the traversal reached: recorded by path, type, size and
// hash, never by content. Content is kept in memory for the judge only.
type Source struct {
	Path       string
	MediaType  string
	Bytes      int64
	SHA256     string
	Depth      int
	LinkedFrom string
	LinkSyntax Syntax
	State      string
	Content    []byte
}

// Result is a traversal: the sources in the order they were reached, how far
// it went, and the context hash of what was read.
type Result struct {
	Sources       []Source
	Files         int
	Depth         int
	Bytes         int64
	ContextSHA256 string
}
