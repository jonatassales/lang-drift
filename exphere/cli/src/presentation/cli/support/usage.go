package support

// UsageError identifies failures caused by invalid CLI input.
// It is shared by the root CLI and individual catalog commands.
type UsageError struct {
	Cause error
}

func (usageError UsageError) Error() string {
	return usageError.Cause.Error()
}

func (usageError UsageError) Unwrap() error {
	return usageError.Cause
}
