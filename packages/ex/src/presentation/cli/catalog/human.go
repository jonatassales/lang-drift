package catalog

import (
	"fmt"
	"sort"
	"strings"

	"github.com/jonatassales/tao/packages/ex/src/application"
)

// The renderers below are presentation only: each writes the same data the
// JSON path emits, shaped for a person. Commands nobody reads interactively
// keep JSON output.

func humanDiagnosis(diagnosis application.Diagnosis) string {
	var report strings.Builder

	report.WriteString("Tools:\n")
	for _, tool := range diagnosis.Tools {
		report.WriteString("  " + humanTool(tool) + "\n")
	}

	report.WriteString("Container engine: " + diagnosis.ContainerEngine + "\n")

	switch {
	case diagnosis.Cluster != nil:
		report.WriteString(fmt.Sprintf("Cluster: %s (context %s, status %s)\n", diagnosis.Cluster.Name, diagnosis.Cluster.Context, diagnosis.Cluster.Status))
	case diagnosis.ClusterNote != "":
		report.WriteString("Cluster: none — " + diagnosis.ClusterNote + "\n")
	}

	report.WriteString("Installation: " + diagnosis.Installation + "\n")

	return report.String()
}

func humanTool(tool application.ToolStatus) string {
	switch tool.Status {
	case application.ToolAvailable:
		line := fmt.Sprintf("%-8s %s %s", "ok", tool.Name, tool.Version)
		if tool.Path != "" {
			line += " (" + tool.Path + ")"
		}
		return line

	case application.ToolMissing:
		line := fmt.Sprintf("%-8s %s", "missing", tool.Name)
		if tool.Install != "" {
			line += " — install: " + tool.Install
		}
		return line

	default:
		return fmt.Sprintf("%-8s %s", tool.Status, tool.Name)
	}
}

func humanCluster(result application.ClusterResult) string {
	var report strings.Builder

	report.WriteString("Cluster " + result.Name + "\n")
	report.WriteString("  Context:    " + result.Context + "\n")
	report.WriteString("  Status:     " + result.Status + "\n")
	report.WriteString("  UID:        " + result.UID + "\n")
	report.WriteString("  Kubeconfig: " + result.Kubeconfig + "\n")
	report.WriteString(fmt.Sprintf("  Namespaces: %d", len(result.NamespaceUIDs)))
	if names := sortedKeys(result.NamespaceUIDs); len(names) > 0 {
		report.WriteString(" (" + strings.Join(names, ", ") + ")")
	}
	report.WriteString("\n")
	report.WriteString(fmt.Sprintf("  Inventory:  version %d, %d resource types\n", result.InventoryVersion, len(result.ResourceTypes)))
	if !result.UpdatedAt.IsZero() {
		report.WriteString("  Updated:    " + result.UpdatedAt.Format("2006-01-02 15:04:05 MST") + "\n")
	}
	if result.Error != "" {
		report.WriteString("  Error:      " + result.Error + "\n")
	}

	return report.String()
}

func sortedKeys(values map[string]string) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	return keys
}
