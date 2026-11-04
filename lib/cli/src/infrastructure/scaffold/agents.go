package scaffold

import (
	"embed"
	"io/fs"
	"path"
	"strings"

	"github.com/jonatassales/tao/lib/cli/src/domain/microruntime"
	"github.com/jonatassales/tao/lib/cli/src/domain/project"
)

// agentTemplates holds the editable TypeScript LangGraph runner shipped with EX.
//
//go:embed assets/agent
var agentTemplates embed.FS

const (
	agentTemplateRoot = "assets/agent"
	templateSuffix    = ".tmpl"
)

// hiddenTemplates are dotfiles stored without their dot, so tooling that
// skips hidden files still embeds them.
var hiddenTemplates = map[string]bool{"gitignore": true, "dockerignore": true}

// Agents publishes agent projects under .langdrift/agents.
type Agents struct{}

// Scaffold copies the embedded templates; nothing is installed or executed.
func (Agents) Scaffold(owner project.Project, name microruntime.Identifier) (string, error) {
	return Publish(owner.Root, path.Join(".langdrift", "agents", name.String()), copyAgentTemplates)
}

func copyAgentTemplates(stage string) error {
	return fs.WalkDir(agentTemplates, agentTemplateRoot, func(source string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil || entry.IsDir() {
			return walkErr
		}

		data, err := agentTemplates.ReadFile(source)
		if err != nil {
			return err
		}

		relative := strings.TrimSuffix(strings.TrimPrefix(source, agentTemplateRoot+"/"), templateSuffix)
		if hiddenTemplates[relative] {
			relative = "." + relative
		}

		return WriteFile(stage, relative, data)
	})
}
