package terraform

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"

	"github.com/exphere/ex/src/domain/microruntime"
)

// Terraform resource types a preparation may manage.
const (
	NamespaceResource = "kubernetes_namespace_v1"
	QuotaResource     = "kubernetes_resource_quota_v1"
)

type plannedChanges struct {
	Changes []struct {
		Type   string `json:"type"`
		Mode   string `json:"mode"`
		Change struct {
			Before  map[string]any `json:"before"`
			After   map[string]any `json:"after"`
			Actions []string       `json:"actions"`
		} `json:"change"`
	} `json:"resource_changes"`
}

// validatePlan refuses a plan that would manage anything but the runtime's
// own namespace and quota, or that would replace them.
func validatePlan(data []byte, namespace, token string) error {
	var plan plannedChanges
	if err := json.Unmarshal(data, &plan); err != nil {
		return fmt.Errorf("decode Terraform plan: %w", err)
	}

	for _, change := range plan.Changes {
		if change.Mode != "managed" {
			return fmt.Errorf("unsupported Terraform data resource %s", change.Type)
		}

		if change.Type != NamespaceResource && change.Type != QuotaResource {
			return fmt.Errorf("unsupported Terraform resource %s", change.Type)
		}

		values := change.Change.After
		if values == nil {
			values = change.Change.Before
		}

		if err := validatePlannedMetadata(change.Type, values, namespace, token); err != nil {
			return err
		}

		replacing := change.Change.Before != nil && change.Change.After != nil
		if replacing && slices.Contains(change.Change.Actions, "delete") {
			return errors.New("replacing the namespace or quota requires cleaning up the runtime first")
		}
	}

	return nil
}

func validatePlannedMetadata(resourceType string, values map[string]any, namespace, token string) error {
	blocks, ok := values["metadata"].([]any)
	if !ok || len(blocks) != 1 {
		return errors.New("planned resource metadata is unknown")
	}

	metadata, ok := blocks[0].(map[string]any)
	if !ok {
		return errors.New("planned resource metadata is invalid")
	}

	labels, _ := metadata["labels"].(map[string]any)
	if labels[microruntime.LabelOwnershipToken] != token {
		return errors.New("planned resource lacks the runtime ownership token")
	}

	if resourceType == NamespaceResource && metadata["name"] != namespace {
		return errors.New("planned namespace differs from the runtime target")
	}

	if resourceType == QuotaResource && metadata["namespace"] != namespace {
		return errors.New("planned quota differs from the runtime target")
	}

	return nil
}
