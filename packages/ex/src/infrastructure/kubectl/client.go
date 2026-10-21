// Package kubectl adapts the consumer-installed kubectl to explicit targets.
// Every call names its kubeconfig and context, so kubectl's ambient context,
// namespace, and kubeconfig are never read.
package kubectl

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/domain/microruntime"
	"github.com/exphere/ex/src/infrastructure/process"
)

const (
	// callTimeout bounds one kubectl invocation; requestTimeout bounds each
	// API request inside it.
	callTimeout    = 30 * time.Second
	requestTimeout = "--request-timeout=15s"
)

// commandRunner runs kubectl directly and propagates cancellation to its
// process group. Run returns after the invocation has settled.
type commandRunner interface {
	Run(context.Context, process.Command) (process.Result, error)
}

// Client implements the application's cluster observation and namespaced
// object ports.
type Client struct {
	Runner commandRunner
}

// run invokes kubectl against access and, when set, one namespace.
func (client Client) run(
	ctx context.Context,
	access application.ClusterAccess,
	namespace string,
	input io.Reader,
	arguments ...string,
) (process.Result, error) {
	callContext, cancel := context.WithTimeout(ctx, callTimeout)
	defer cancel()

	target := []string{"--kubeconfig", access.Kubeconfig, "--context", access.Context, requestTimeout}
	if namespace != "" {
		target = append(target, "--namespace", namespace)
	}

	result, err := client.Runner.Run(callContext, process.Command{
		Name:  "kubectl",
		Args:  append(target, arguments...),
		Stdin: input,
	})
	if err != nil {
		return result, process.Explain("kubectl "+arguments[0], result, err)
	}

	return result, nil
}

// decodeItems reads the items of a kubectl list document.
func decodeItems(resource, document string) ([]microruntime.Object, error) {
	var list struct {
		Items []microruntime.Object `json:"items"`
	}

	if err := json.Unmarshal([]byte(document), &list); err != nil {
		return nil, fmt.Errorf("decode %s: %w", resource, err)
	}

	return list.Items, nil
}

// isEvent recognizes events, which describe objects but are not owned state.
func isEvent(resource string) bool {
	return resource == "events" || resource == "events.events.k8s.io"
}
