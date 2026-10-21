package kubectl

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/domain/microruntime"
	"github.com/exphere/ex/src/infrastructure/process"
)

// Logs streams the selected workload's native logs. Follow is bounded by the
// caller's cancellation; a finite read additionally limits output to one MiB.
// The ordinary per-call API timeout must not shorten an explicit live stream.
func (client Client) Logs(ctx context.Context, target application.Target, resource microruntime.Resource, follow bool, output, diagnostics io.Writer) error {
	arguments := append(streamTarget(target), "logs", strings.ToLower(resource.Kind)+"/"+resource.Name,
		"--all-containers=true", "--tail=500", "--timestamps=true")
	if follow {
		arguments = append(arguments, "--follow=true")
	} else {
		arguments = append(arguments, "--limit-bytes=1048576")
	}
	_, err := client.Runner.Run(ctx, process.Command{
		Name: "kubectl", Args: arguments, Stdout: output, Stderr: diagnostics,
	})

	return err
}

// Forward owns a foreground localhost-only tunnel and waits for cancellation
// or native termination. Both native streams are diagnostics, not results.
func (client Client) Forward(ctx context.Context, target application.Target, resource microruntime.Resource, localPort, remotePort int, diagnostics io.Writer) error {
	arguments := append(streamTarget(target), "port-forward", "service/"+resource.Name,
		fmt.Sprintf("%d:%d", localPort, remotePort), "--address=127.0.0.1")
	_, err := client.Runner.Run(ctx, process.Command{
		Name: "kubectl", Args: arguments, Stdout: diagnostics, Stderr: diagnostics,
	})

	return err
}

func streamTarget(target application.Target) []string {
	return []string{"--kubeconfig", target.Kubeconfig, "--context", target.Context, "--namespace", target.Namespace}
}
