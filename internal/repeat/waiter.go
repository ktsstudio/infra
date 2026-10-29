package repeat

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v5"
)

// Waiter tracks the number of failures and waits for a duration of time based
// on the number of failures.
type Waiter struct {
	backOff    backoff.BackOff
	started    time.Time
	maxElapsed time.Duration
}

func NewWaiter(backOff backoff.BackOff) *Waiter {
	return NewWaiterWithMaxElapsed(backOff, 0)
}

// NewWaiterWithMaxElapsed stops retries after maxElapsed; zero retries indefinitely.
func NewWaiterWithMaxElapsed(backOff backoff.BackOff, maxElapsed time.Duration) *Waiter {
	// Initialize the backoff here, so we don't have to remember to do it in
	// other places.
	if e, ok := backOff.(*backoff.ExponentialBackOff); ok {
		if e.MaxInterval == 0 {
			panic("exponential backoff requires a maximum interval")
		}
		e.Reset()
	}
	return &Waiter{backOff: backOff, started: time.Now(), maxElapsed: maxElapsed}
}

// Reset sets the number of failures to 0, so that the next call to Wait sleeps
// for the minimum duration.
func (b *Waiter) Reset() {
	b.backOff.Reset()
	b.started = time.Now()
}

// Wait blocks for the duration of delay calculated by BackOff, or until the
// context is done.
// Returns an error when the context is done, or the retry limit is reached.
// Otherwise, returns nil when the timer waited the full duration.
func (b Waiter) Wait(ctx context.Context) error {
	delay := b.backOff.NextBackOff()
	if delay == backoff.Stop || (b.maxElapsed > 0 && time.Since(b.started)+delay > b.maxElapsed) {
		return fmt.Errorf("retry limit exceeded")
	}
	timer := time.NewTimer(delay)
	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		timer.Stop()
		return ctx.Err()
	}
}
