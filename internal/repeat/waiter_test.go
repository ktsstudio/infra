package repeat

import (
	"context"
	"testing"
	"time"

	"github.com/cenkalti/backoff/v5"
	"gotest.tools/v3/assert"
)

func TestWaiter_Reset(t *testing.T) {
	backOff := &backoff.ExponentialBackOff{
		InitialInterval: 1,
		MaxInterval:     100,
		Multiplier:      2,
	}
	w := NewWaiter(backOff)
	assert.Equal(t, backOff.NextBackOff(), time.Duration(1))
	assert.Equal(t, backOff.NextBackOff(), time.Duration(2))
	w.Reset()
	assert.Equal(t, backOff.NextBackOff(), time.Duration(1))
}

func TestWaiter_Wait(t *testing.T) {
	t.Run("error when context done", func(t *testing.T) {
		w := NewWaiter(backoff.NewConstantBackOff(time.Minute))
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		assert.ErrorIs(t, w.Wait(ctx), context.Canceled)
	})
	t.Run("error when backoff returns stop", func(t *testing.T) {
		w := NewWaiter(&backoff.StopBackOff{})
		assert.Error(t, w.Wait(context.Background()), "retry limit exceeded")

	})
	t.Run("no error on timer tick", func(t *testing.T) {
		w := NewWaiter(backoff.NewConstantBackOff(time.Millisecond))
		assert.NilError(t, w.Wait(context.Background()))
	})
}

func TestWaiterMaxElapsed(t *testing.T) {
	waiter := NewWaiterWithMaxElapsed(backoff.NewConstantBackOff(time.Nanosecond), time.Minute)
	waiter.started = time.Now().Add(-2 * time.Minute)
	assert.Error(t, waiter.Wait(context.Background()), "retry limit exceeded")
	waiter.Reset()
	assert.NilError(t, waiter.Wait(context.Background()))
}
