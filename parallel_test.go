package statstest

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestParallelComputeEachContextRunsEveryIndex(t *testing.T) {
	const n = 50
	var seen [n]int32
	err := parallelComputeEachContext(context.Background(), n, func(i int) {
		atomic.AddInt32(&seen[i], 1)
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i, v := range seen {
		if v != 1 {
			t.Fatalf("index %d ran %d times, want exactly 1", i, v)
		}
	}
}

func TestParallelComputeEachContextSequentialBelowThreshold(t *testing.T) {
	const n = parallelThreshold // at or below threshold must run sequentially
	var seen [n]int32
	err := parallelComputeEachContext(context.Background(), n, func(i int) {
		atomic.AddInt32(&seen[i], 1)
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i, v := range seen {
		if v != 1 {
			t.Fatalf("index %d ran %d times, want exactly 1", i, v)
		}
	}
}

func TestParallelComputeEachContextStopsOnCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // already canceled before dispatch starts

	var calls int32
	err := parallelComputeEachContext(ctx, 100, func(i int) {
		atomic.AddInt32(&calls, 1)
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if got := atomic.LoadInt32(&calls); got >= 100 {
		t.Fatalf("expected cancellation to skip at least some of the 100 calls, got %d", got)
	}
}

func TestParallelComputeEachContextRespectsDeadlineMidRun(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	var calls int32
	err := parallelComputeEachContext(ctx, 1000, func(i int) {
		atomic.AddInt32(&calls, 1)
		time.Sleep(time.Millisecond)
	})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected context.DeadlineExceeded, got %v", err)
	}
	if got := atomic.LoadInt32(&calls); got >= 1000 {
		t.Fatalf("expected the deadline to cut the run short, got all %d calls", got)
	}
}

func TestTukeyHSDContextCancellationBeforeCall(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := TukeyHSDContext(ctx, 0.95, []float64{1, 2, 3}, []float64{4, 5, 6})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestTukeyHSDContextMatchesTukeyHSD(t *testing.T) {
	g1 := []float64{8.1, 8.3, 7.9, 8.0, 8.2}
	g2 := []float64{8.8, 9.0, 8.7, 8.9, 9.1}
	g3 := []float64{7.5, 7.6, 7.4, 7.7, 7.5}

	want, err := TukeyHSD(0.95, g1, g2, g3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := TukeyHSDContext(context.Background(), 0.95, g1, g2, g3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got.Comparisons) != len(want.Comparisons) {
		t.Fatalf("got %d comparisons, want %d", len(got.Comparisons), len(want.Comparisons))
	}
	for i := range want.Comparisons {
		if got.Comparisons[i] != want.Comparisons[i] {
			t.Fatalf("comparison %d: got %+v, want %+v", i, got.Comparisons[i], want.Comparisons[i])
		}
	}
}
