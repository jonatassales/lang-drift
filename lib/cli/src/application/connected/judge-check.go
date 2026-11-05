package connected

import "context"

// JudgeCheck is the connection check against a console's judge route: it
// accepts the deployment secret, or the connection is refused naming the
// console.
type JudgeCheck interface {
	Check(ctx context.Context, consoleURL string, secret string) error
}
