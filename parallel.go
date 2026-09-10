package statstest

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
)

// parallelThreshold is the minimum number of independent units of work
// before parallelComputeEachContext bothers spinning up goroutines. Below
// this, goroutine scheduling overhead outweighs any benefit.
const parallelThreshold = 4

// parallelComputeEachContext calls f(i) once for every i in [0, n),
// distributing the calls across up to GOMAXPROCS goroutines when there is
// enough independent work to be worth it, and running sequentially
// otherwise. Passing context.Background() makes cancellation a no-op, so
// every index always runs and the returned error is always nil.
//
// f must be safe to call concurrently and must only write to state owned
// by index i (e.g. results[i]) — parallelComputeEachContext does no
// synchronization beyond ensuring every index is dispatched at most once
// and waiting for all dispatched calls to finish before returning.
//
// It stops dispatching new work as soon as ctx is done and returns
// ctx.Err() in that case (nil once every index has actually run). Work
// already in-flight when ctx is canceled is not interrupted mid-call —
// only the dispatch of further indices stops — so f itself is never asked
// to observe cancellation.
func parallelComputeEachContext(ctx context.Context, n int, f func(i int)) error {
	if n <= parallelThreshold || runtime.GOMAXPROCS(0) <= 1 {
		for i := 0; i < n; i++ {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
			f(i)
		}
		return nil
	}

	workers := runtime.GOMAXPROCS(0)
	if workers > n {
		workers = n
	}

	var next int64 = -1
	var wg sync.WaitGroup
	wg.Add(workers)
	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				default:
				}
				i := int(atomic.AddInt64(&next, 1))
				if i >= n {
					return
				}
				f(i)
			}
		}()
	}
	wg.Wait()

	return ctx.Err()
}
