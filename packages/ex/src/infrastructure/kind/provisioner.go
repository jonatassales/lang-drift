// Package kind adapts the consumer-installed kind executable. It never
// installs tools, builds images, or changes the user's ambient kubeconfig.
package kind

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/exphere/ex/src/infrastructure/process"
)

const (
	// contextPrefix is how kind names a cluster's kubeconfig context.
	contextPrefix = "kind-"

	// readinessWait bounds how long kind waits for the control plane.
	readinessWait = "120s"

	credentialsMode = 0o600
)

// commandRunner runs kind or Docker directly and propagates cancellation to
// its process group. Run returns after the invocation has settled.
type commandRunner interface {
	Run(context.Context, process.Command) (process.Result, error)
}

// Provisioner creates, deletes, and loads images into local kind clusters.
type Provisioner struct {
	Runner commandRunner
}

// Clusters lists the names of existing kind clusters.
func (provisioner Provisioner) Clusters(ctx context.Context) ([]string, error) {
	result, err := provisioner.kind(ctx, "list kind clusters", "get", "clusters")
	if err != nil {
		return nil, err
	}

	return strings.Fields(result.Stdout), nil
}

// Create writes the cluster's credentials only to kubeconfig, a private
// project file, and returns the cluster's context name. A partially created
// cluster may remain when creation fails.
func (provisioner Provisioner) Create(ctx context.Context, name, kubeconfig string) (string, error) {
	if _, err := provisioner.kind(ctx, "create kind cluster",
		"create", "cluster", "--name", name, "--kubeconfig", kubeconfig, "--wait", readinessWait); err != nil {
		return "", err
	}

	if err := os.Chmod(kubeconfig, credentialsMode); err != nil {
		return "", err
	}

	return contextPrefix + name, nil
}

// Delete deletes the named cluster. Deciding that it may be deleted is the
// caller's responsibility; this adapter never infers a target.
func (provisioner Provisioner) Delete(ctx context.Context, name string) error {
	_, err := provisioner.kind(ctx, "delete kind cluster", "delete", "cluster", "--name", name)

	return err
}

// LoadImage imports an already-built local Docker image into the cluster.
func (provisioner Provisioner) LoadImage(ctx context.Context, clusterName, image string) error {
	if err := requireOperand(image); err != nil {
		return err
	}

	_, err := provisioner.kind(ctx, "load image", "load", "docker-image", image, "--name", clusterName)

	return err
}

// RequireImages fails unless every image is present on every cluster node.
func (provisioner Provisioner) RequireImages(ctx context.Context, clusterName string, images []string) error {
	nodes, err := provisioner.nodes(ctx, clusterName)
	if err != nil {
		return err
	}

	for _, image := range images {
		if err := requireOperand(image); err != nil {
			return err
		}

		for _, node := range nodes {
			probe := process.Command{Name: "docker", Args: []string{"exec", node, "crictl", "inspecti", image}}

			if result, err := provisioner.Runner.Run(ctx, probe); err != nil {
				return process.Explain(fmt.Sprintf("image %q is not loaded on node %s", image, node), result, err)
			}
		}
	}

	return nil
}

func (provisioner Provisioner) nodes(ctx context.Context, clusterName string) ([]string, error) {
	result, err := provisioner.kind(ctx, "list kind nodes", "get", "nodes", "--name", clusterName)
	if err != nil {
		return nil, err
	}

	nodes := strings.Fields(result.Stdout)
	if len(nodes) == 0 {
		return nil, fmt.Errorf("kind cluster %s has no nodes", clusterName)
	}

	return nodes, nil
}

func (provisioner Provisioner) kind(ctx context.Context, action string, arguments ...string) (process.Result, error) {
	result, err := provisioner.Runner.Run(ctx, process.Command{Name: "kind", Args: arguments})
	if err != nil {
		return result, process.Explain(action, result, err)
	}

	return result, nil
}

// requireOperand refuses a value a native tool would parse as an option.
func requireOperand(value string) error {
	if strings.HasPrefix(value, "-") {
		return errors.New("image reference must not start with '-'")
	}

	return nil
}
