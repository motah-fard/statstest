package statstest

import (
	"runtime"
	"sync"
	"sync/atomic"
)

// parallelThreshold is the minimum number of independent units of work
// before parallelComputeEach bothers spinning up goroutines. Below this,
// goroutine scheduling overhead outweighs any benefit.
const parallelThreshold = 4

// parallelComputeEach calls f(i) once for every i in [0, n), distributing
// the calls across up to GOMAXPROCS goroutines when there is enough
// independent work to be worth it, and running sequentially otherwise.
//
// f must be safe to call concurrently and must only write to state owned
// by index i (e.g. results[i]) — parallelComputeEach does no
// synchronization beyond ensuring every index is dispatched exactly once
// and waiting for all of them to finish before returning.
func parallelComputeEach(n int, f func(i int)) {
	if n <= parallelThreshold || runtime.GOMAXPROCS(0) <= 1 {
		for i := 0; i < n; i++ {
			f(i)
		}
		return
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
				i := int(atomic.AddInt64(&next, 1))
				if i >= n {
					return
				}
				f(i)
			}
		}()
	}
	wg.Wait()
}
