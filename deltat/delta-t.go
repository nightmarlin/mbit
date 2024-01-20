// Package deltat ("delta T") implements simple "elapsed-time" loop helpers.
// The "delta" can be any number type, and the calculation provided arbitrarily.
// deltat also provides basic implementations of this
package deltat

import (
	"runtime"
	"time"

	"golang.org/x/exp/constraints"
)

// A WallClock returns the current system time.
type WallClock func() time.Time

// The DefaultClock is a WallClock that returns the result of [time.Now]
func DefaultClock() WallClock { return time.Now }

type (
	// A Numeric is any float or int type.
	Numeric interface {
		constraints.Float | constraints.Integer
	}

	// A DeltaCalculator returns the elapsed D between prev and now
	DeltaCalculator[D Numeric] func(prev, now time.Time) (delta D)

	// A DeltaFn performs some action based on the delta provided. It returns
	// whether it should continue into the next iteration.
	DeltaFn[D Numeric] func(now time.Time, delta D) (cont bool)
)

// The Loop calculates the elapsed D between prev and now using the provided
// DeltaCalculator, then passes that value to the provided DeltaFn. It exits
// when the DeltaFn returns `false`.
func Loop[D constraints.Integer | constraints.Float](
	clock WallClock,
	calcDelta DeltaCalculator[D],
	do DeltaFn[D],
) {
	prev := clock()
	for {
		now := clock()
		elapsed := calcDelta(prev, now)
		prev = now

		if !do(now, elapsed) {
			return
		}
		runtime.Gosched() // allow other routines a chance to run
	}
}

// ByTick is a DeltaCalculator that returns how many ticks have elapsed since
// the last iteration. tps sets this value.
func ByTick(tps float32) DeltaCalculator[float32] {
	return func(prev, now time.Time) (delta float32) {
		return (float32(now.Sub(prev)) / float32(time.Second)) * tps
	}
}

// ByElapsed is a DeltaCalculator that simply returns the elapsed
// [time.Duration] since the previous iteration.
func ByElapsed() DeltaCalculator[time.Duration] {
	return func(prev, now time.Time) (delta time.Duration) { return now.Sub(prev) }
}
