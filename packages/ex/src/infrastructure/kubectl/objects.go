package kubectl

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/domain/microruntime"
)

const fieldManager = "--field-manager=exphere"

// plurals maps the kinds EX deletes to their API resource names.
var plurals = map[string]string{
	"Deployment": "deployments",
	"Service":    "services",
	"Job":        "jobs",
}

// Get reads one object; an absent object is nil, an API failure an error.
func (client Client) Get(ctx context.Context, target application.Target, kind, name string) (microruntime.Object, error) {
	result, err := client.run(ctx, target.ClusterAccess, target.Namespace, nil,
		"get", kind, name, "--ignore-not-found=true", "-o", "json")
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(result.Stdout) == "" {
		return nil, nil
	}

	var object microruntime.Object
	if err := json.Unmarshal([]byte(result.Stdout), &object); err != nil {
		return nil, fmt.Errorf("decode %s/%s: %w", kind, name, err)
	}

	return object, nil
}

// List reads every object of one resource type in the target namespace.
func (client Client) List(ctx context.Context, target application.Target, resource string) ([]microruntime.Object, error) {
	result, err := client.run(ctx, target.ClusterAccess, target.Namespace, nil, "get", resource, "-o", "json")
	if err != nil {
		return nil, err
	}

	return decodeItems(resource, result.Stdout)
}

// DryRun validates object against the live API without persisting it.
func (client Client) DryRun(ctx context.Context, target application.Target, object microruntime.Object) error {
	_, err := client.submit(ctx, target, object,
		"apply", "--server-side", fieldManager, "--dry-run=server", "-f", "-", "-o", "json")

	return err
}

// Create creates object. A live name collision fails rather than adopting
// an object someone else created meanwhile.
func (client Client) Create(ctx context.Context, target application.Target, object microruntime.Object) (microruntime.Object, error) {
	return client.submit(ctx, target, object, "create", fieldManager, "-f", "-", "-o", "json")
}

func (client Client) submit(ctx context.Context, target application.Target, object microruntime.Object, arguments ...string) (microruntime.Object, error) {
	document, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}

	result, err := client.run(ctx, target.ClusterAccess, target.Namespace, bytes.NewReader(document), arguments...)
	if err != nil {
		return nil, err
	}

	var returned microruntime.Object
	if err := json.Unmarshal([]byte(result.Stdout), &returned); err != nil {
		return nil, fmt.Errorf("decode %s result: %w", arguments[0], err)
	}

	return returned, nil
}

// Delete requests deletion of exactly the recorded object.
func (client Client) Delete(ctx context.Context, target application.Target, resource microruntime.Resource) error {
	plural, supported := plurals[resource.Kind]
	if !supported {
		return fmt.Errorf("deleting %s objects is unsupported", resource.Kind)
	}

	base := "/api/v1"
	if resource.APIVersion != "v1" {
		base = "/apis/" + resource.APIVersion
	}

	path := base + "/namespaces/" + resource.Namespace + "/" + plural + "/" + resource.Name

	return client.deleteExactly(ctx, target, path, resource.UID)
}

// DeleteNamespace requests deletion of exactly the recorded target namespace.
func (client Client) DeleteNamespace(ctx context.Context, target application.Target, uid string) error {
	if uid == "" {
		return errors.New("namespace deletion requires an observed UID")
	}

	return client.deleteExactly(ctx, target, "/api/v1/namespaces/"+target.Namespace, uid)
}

// deleteExactly deletes in the foreground only the object whose UID matches,
// so a replacement under the same name is never deleted by mistake.
func (client Client) deleteExactly(ctx context.Context, target application.Target, path, uid string) error {
	options, err := json.Marshal(map[string]any{
		"apiVersion":        "v1",
		"kind":              "DeleteOptions",
		"preconditions":     map[string]string{"uid": uid},
		"propagationPolicy": "Foreground",
	})
	if err != nil {
		return err
	}

	_, err = client.run(ctx, target.ClusterAccess, target.Namespace, bytes.NewReader(options),
		"delete", "--raw", path, "-f", "-")

	return err
}

// NamespaceInventory lists every object of every listable and deletable
// resource type in the target namespace, except events. Views that cannot be
// deleted, such as pod metrics, hold no state and carry no UID, so they are
// excluded as in the cluster inventory. A partial inventory is an error.
func (client Client) NamespaceInventory(ctx context.Context, target application.Target) ([]microruntime.Object, error) {
	result, err := client.run(ctx, target.ClusterAccess, target.Namespace, nil,
		"api-resources", "--verbs=list,delete", "--namespaced=true", "-o", "name")
	if err != nil {
		return nil, err
	}

	var inventory []microruntime.Object

	for _, resource := range strings.Fields(result.Stdout) {
		if isEvent(resource) {
			continue
		}

		objects, err := client.List(ctx, target, resource)
		if err != nil {
			return nil, fmt.Errorf("incomplete inventory (%s): %w", resource, err)
		}

		inventory = append(inventory, objects...)
	}

	return inventory, nil
}

// PodMetrics reads pod metrics from the optional metrics API.
func (client Client) PodMetrics(ctx context.Context, target application.Target) ([]microruntime.Object, error) {
	result, err := client.run(ctx, target.ClusterAccess, target.Namespace, nil,
		"get", "--raw", "/apis/metrics.k8s.io/v1beta1/namespaces/"+target.Namespace+"/pods")
	if err != nil {
		return nil, err
	}

	return decodeItems("pod metrics", result.Stdout)
}
