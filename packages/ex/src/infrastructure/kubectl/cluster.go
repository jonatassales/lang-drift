package kubectl

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/domain/cluster"
)

// inventoryProjection prints one line per object: namespace, name, UID, and
// the comma-terminated UIDs of its controller owners. No payload is read.
const inventoryProjection = `{{range .items}}` +
	`{{if .metadata.namespace}}{{.metadata.namespace}}{{end}}{{"\t"}}` +
	`{{.metadata.name}}{{"\t"}}{{.metadata.uid}}{{"\t"}}` +
	`{{range .metadata.ownerReferences}}{{if .controller}}{{.uid}}{{","}}{{end}}{{end}}{{"\n"}}` +
	`{{end}}`

// missingValue is what a Go template prints for an absent field.
const missingValue = "<no value>"

// Namespaces maps every namespace name to its UID.
func (client Client) Namespaces(ctx context.Context, access application.ClusterAccess) (map[string]string, error) {
	result, err := client.run(ctx, access, "", nil, "get", "namespaces", "-o", "json")
	if err != nil {
		return nil, err
	}

	items, err := decodeItems("namespaces", result.Stdout)
	if err != nil {
		return nil, err
	}

	namespaces := make(map[string]string, len(items))

	for _, item := range items {
		if item.Name() == "" || item.UID() == "" {
			return nil, errors.New("namespace inventory lacks an observed identity")
		}

		namespaces[item.Name()] = item.UID()
	}

	return namespaces, nil
}

// ClusterInventory observes every persisted object of every listable and
// deletable resource type, namespaced or not. Virtual API views are not
// listable and deletable, so they are excluded; events are skipped.
func (client Client) ClusterInventory(ctx context.Context, access application.ClusterAccess) (cluster.Inventory, error) {
	inventory := cluster.Inventory{Objects: map[string]cluster.ObjectIdentity{}}

	for _, namespaced := range []bool{true, false} {
		resources, err := client.resourceTypes(ctx, access, namespaced)
		if err != nil {
			return cluster.Inventory{}, err
		}

		for _, resource := range resources {
			if isEvent(resource) {
				continue
			}

			inventory.ResourceTypes = append(inventory.ResourceTypes, strconv.FormatBool(namespaced)+"/"+resource)

			if err := client.collect(ctx, access, resource, namespaced, inventory.Objects); err != nil {
				return cluster.Inventory{}, fmt.Errorf("incomplete cluster inventory (%s): %w", resource, err)
			}
		}
	}

	if len(inventory.Objects) == 0 {
		return cluster.Inventory{}, errors.New("an empty object inventory cannot establish cluster ownership")
	}

	return inventory, nil
}

func (client Client) resourceTypes(ctx context.Context, access application.ClusterAccess, namespaced bool) ([]string, error) {
	result, err := client.run(ctx, access, "", nil,
		"api-resources", "--namespaced="+strconv.FormatBool(namespaced), "--verbs=list,delete", "-o", "name")
	if err != nil {
		return nil, err
	}

	resources := strings.Fields(result.Stdout)
	if len(resources) == 0 {
		return nil, errors.New("empty API discovery cannot establish complete cluster ownership")
	}

	sort.Strings(resources)

	return resources, nil
}

func (client Client) collect(
	ctx context.Context,
	access application.ClusterAccess,
	resource string,
	namespaced bool,
	objects map[string]cluster.ObjectIdentity,
) error {
	arguments := []string{"get", resource, "-o", "go-template=" + inventoryProjection}
	if namespaced {
		arguments = append(arguments, "--all-namespaces")
	}

	result, err := client.run(ctx, access, "", nil, arguments...)
	if err != nil {
		return err
	}

	for _, line := range strings.Split(strings.TrimSuffix(result.Stdout, "\n"), "\n") {
		if line == "" {
			continue
		}

		fields := strings.Split(line, "\t")
		if len(fields) != 4 || fields[1] == "" || fields[2] == "" || fields[2] == missingValue ||
			(namespaced && fields[0] == "") {
			return fmt.Errorf("resource %s lacks an observable persisted identity", resource)
		}

		key := resource + "/" + fields[0] + "/" + fields[1]
		if _, duplicate := objects[key]; duplicate {
			return fmt.Errorf("duplicate identity in cluster inventory: %s", key)
		}

		objects[key] = cluster.ObjectIdentity{UID: fields[2], Controllers: controllerUIDs(fields[3])}
	}

	return nil
}

func controllerUIDs(field string) []string {
	var controllers []string

	for _, uid := range strings.Split(strings.TrimSuffix(field, ","), ",") {
		if uid != "" && uid != missingValue {
			controllers = append(controllers, uid)
		}
	}

	return controllers
}
