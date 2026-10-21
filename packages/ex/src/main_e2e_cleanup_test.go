//go:build e2e

package main

import (
	"context"
	"slices"
	"strings"
	"time"

	domainproject "github.com/jonatassales/tao/packages/ex/src/domain/project"
)

// foreignLabelKey marks the namespaces the run creates to stand for objects
// EX did not create; its value is the run nonce.
const foreignLabelKey = "langdrift.e2e/foreign"

// foreignSelector selects exactly the foreign namespaces this run created.
func (h *harness) foreignSelector() string {
	return foreignLabelKey + "=" + h.nonce
}

// cleanup removes only what the run identified, through the CLI's own
// ownership checks, and reports what it cannot remove safely. It never
// deletes a cluster or namespace by name alone. Whether something is left
// behind is decided by observing it; a failed observation counts as left
// behind, because absence could not be confirmed.
func (h *harness) cleanup() {
	ctx, cancel := context.WithTimeout(context.Background(), cleanupBudget)
	defer cancel()

	clusterReachable := h.cluster.Kubeconfig != "" && h.cluster.Context != ""

	if clusterReachable && h.foreignCreated {
		started := time.Now()
		result := h.deleteForeignNamespaces(ctx)
		h.record(&h.report.Cleanup, "remove labeled foreign namespaces",
			[]string{"kubectl", "delete", "namespaces", "-l", h.foreignSelector()}, result, time.Since(started), nil)

		remaining := h.kubectl(ctx, nil, "get", "namespaces", "-l", h.foreignSelector(), "-o", "name")
		h.requireObservedAbsent(remaining, "namespaces labeled "+h.foreignSelector())
	}

	if h.upAttempted {
		h.runCLI(ctx, &h.report.Cleanup, "runtime down --recover", "", h.finiteCommand("runtime", "down", runtimeID, "--recover"))

		if clusterReachable {
			remaining := h.kubectl(ctx, nil, "get", "namespace", h.namespace(), "--ignore-not-found=true", "-o", "name")
			h.requireObservedAbsent(remaining, "namespace "+h.namespace())
		}
	}

	if h.clusterAttempted {
		h.cleanupCluster(ctx)
	}

	for tag, id := range h.images {
		h.removeImage(ctx, tag, id)
	}

	h.finish()
}

// requireObservedAbsent reports a resource unless a successful observation
// found nothing.
func (h *harness) requireObservedAbsent(observed outcome, resource string) {
	switch {
	case observed.exit != 0:
		h.reportLeftBehind(resource, "could not confirm removal: "+firstLine(observed.stderr))
	case strings.TrimSpace(observed.stdout) != "":
		h.reportLeftBehind(resource, "still present after cleanup")
	}
}

// deleteForeignNamespaces deletes only the namespaces carrying this run's
// label and waits until they are gone.
func (h *harness) deleteForeignNamespaces(ctx context.Context) outcome {
	return h.kubectl(ctx, nil, "delete", "namespaces", "-l", h.foreignSelector(),
		"--ignore-not-found=true", "--wait=true", "--timeout=180s")
}

func (h *harness) cleanupCluster(ctx context.Context) {
	name := h.cluster.Name
	if name == "" {
		name = h.expectedClusterName()
	}

	h.runCLI(ctx, &h.report.Cleanup, "cluster delete", "", h.finiteCommand("cluster", "delete"))

	if name == "" {
		h.reportLeftBehind("kind cluster of this run", "its name is unknown; check `kind get clusters` for an ex- cluster")
		return
	}

	listed := h.native(ctx, "kind", nil, "get", "clusters")

	switch {
	case listed.exit != 0:
		h.reportLeftBehind("kind cluster "+name, "could not confirm removal: "+firstLine(listed.stderr))
	case slices.Contains(strings.Fields(listed.stdout), name):
		h.reportLeftBehind("kind cluster "+name, "still present; never deleted by name without ownership reconciliation")
	}
}

// expectedClusterName derives the cluster name from the project identity
// with the domain's own rule, for a creation that failed before reporting it.
func (h *harness) expectedClusterName() string {
	id, err := domainproject.ParseID(h.projectID)
	if err != nil {
		return ""
	}

	return domainproject.Project{ID: id}.ClusterName()
}

// removeImage removes a tag only while it still names the recorded image,
// then confirms that it is gone.
func (h *harness) removeImage(ctx context.Context, tag, id string) {
	started := time.Now()
	listed := h.imageID(ctx, tag)

	switch observed := strings.TrimSpace(listed.stdout); {
	case listed.exit != 0:
		h.reportLeftBehind("image "+tag, "could not confirm removal: "+firstLine(listed.stderr))
		return

	case observed == "":
		h.record(&h.report.Cleanup, "image "+tag+" already absent", []string{"docker", "image", "ls", tag},
			listed, time.Since(started), nil)
		return

	case observed != id:
		h.reportLeftBehind("image "+tag, "the tag now names a different image; not removed")
		return
	}

	result := h.native(ctx, "docker", nil, "image", "rm", tag)
	h.record(&h.report.Cleanup, "remove image "+tag, []string{"docker", "image", "rm", tag}, result, time.Since(started), nil)

	h.requireObservedAbsent(h.imageID(ctx, tag), "image "+tag)
}

// reportLeftBehind records a resource cleanup could not remove safely.
func (h *harness) reportLeftBehind(resource, reason string) {
	h.report.LeftBehind = append(h.report.LeftBehind, leftBehind{Resource: resource, Reason: sanitize(reason, h.workdir)})
}
