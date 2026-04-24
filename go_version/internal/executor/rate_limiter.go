package executor

import (
	"context"

	"golang.org/x/time/rate"
)

type RateLimitedExecutor struct {
	limiter  *rate.Limiter
	delegate Executor
}

func NewRateLimitedExecutor(maxPerSec int) *RateLimitedExecutor {
	if maxPerSec <= 0 {
		maxPerSec = 1
	}
	return &RateLimitedExecutor{
		limiter:  rate.NewLimiter(rate.Limit(maxPerSec), maxPerSec),
		delegate: NewShellExecutor(),
	}
}

func (r *RateLimitedExecutor) Execute(ctx context.Context, command string, opts ExecOptions) *ExecResult {
	if err := r.limiter.Wait(ctx); err != nil {
		return &ExecResult{
			Command:  command,
			ExitCode: -1,
			Stderr:   "rate limit: " + err.Error(),
		}
	}
	return r.delegate.Execute(ctx, command, opts)
}
