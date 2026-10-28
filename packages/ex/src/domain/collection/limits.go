package collection

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"path"
	"sort"
	"strings"
)

// Limits bound the traversal.
type Limits struct {
	MaxDepth int
	MaxFiles int
	MaxBytes int64
}

// DefaultLimits are the limits a plan takes unless the owner says otherwise.
func DefaultLimits() Limits {
	return Limits{MaxDepth: 3, MaxFiles: 64, MaxBytes: 4 << 20}
}

// The states a source may be in.
const (
	Read       = "read"
	Binary     = "binary"
	Unresolved = "unresolved"
	Refused    = "refused"
)

// What the limits say about what they refuse.
const RefusedWord = "refused"

// Media types the records name.
const (
	MarkdownType = "text/markdown"
	JSONType     = "application/json"
	PlainType    = "text/plain"
	BinaryType   = "application/octet-stream"
)

// The private state of the folder, which no plan may read.
const privateState = ".langdrift/agent/state"

var secretSuffixes = []string{".pem", ".key", ".p12", ".pfx", ".jks", ".keystore"}

// IsSecret reports whether a path names something a plan never reads: a
// dotenv file, a key, a credential, or the folder's private state.
func IsSecret(relativePath string) bool {
	if relativePath == privateState || strings.HasPrefix(relativePath, privateState+"/") {
		return true
	}

	name := strings.ToLower(path.Base(relativePath))
	if name == ".env" || strings.HasPrefix(name, ".env.") || strings.HasPrefix(name, "id_rsa") || strings.HasPrefix(name, "id_ed25519") {
		return true
	}

	if strings.Contains(name, "credential") || strings.Contains(name, "secret") {
		return true
	}

	for _, suffix := range secretSuffixes {
		if strings.HasSuffix(name, suffix) {
			return true
		}
	}

	return false
}

// MediaType is the type the records name a text file by, from its name.
func MediaType(relativePath string) string {
	switch strings.ToLower(path.Ext(relativePath)) {
	case ".md", ".markdown":
		return MarkdownType
	case ".json":
		return JSONType
	}

	return PlainType
}

// IsBinary reports whether content is not text: a NUL byte in its first
// 8 KiB decides it.
func IsBinary(content []byte) bool {
	sample := content
	if len(sample) > 8<<10 {
		sample = sample[:8<<10]
	}

	return bytes.IndexByte(sample, 0) >= 0
}

// ContextRule is how the context hash is taken, as the records say it.
const ContextRule = "sha256 over the sorted path<TAB>sha256 lines of the sources read; byte hashes until canonical hashing lands"

// ContextHash hashes what the plan read: one line per source, its path and
// hash separated by a tab, sorted, each ended by a newline.
func ContextHash(hashes map[string]string) string {
	lines := make([]string, 0, len(hashes))
	for relativePath, hash := range hashes {
		lines = append(lines, relativePath+"\t"+hash+"\n")
	}

	sort.Strings(lines)

	return SHA256([]byte(strings.Join(lines, "")))
}

// SHA256 is the hex digest of content.
func SHA256(content []byte) string {
	digest := sha256.Sum256(content)

	return hex.EncodeToString(digest[:])
}
