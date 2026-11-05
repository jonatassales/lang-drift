package plan

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// promptHash is the SHA-256 of the system text and the three prompts, as the
// plan records it beside PromptVersion: a prompt edited without a new
// version shows as a new hash.
func promptHash() string {
	prompts := strings.Join([]string{judgeSystem, sourcesPrompt, episodesPrompt, suggestionPrompt}, "\n\x00\n")
	digest := sha256.Sum256([]byte(prompts))

	return "sha256:" + hex.EncodeToString(digest[:])
}
