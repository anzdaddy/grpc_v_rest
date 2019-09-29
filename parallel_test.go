package main

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
)

func parallelBenchmark(
	b *testing.B,
	parallelism int,
	worker func(ctx context.Context, work func() bool) error,
) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errch := make(chan error, parallelism)
	panicch := make(chan interface{}, parallelism)

	var totalSuccesses uint64 = 0

	var wg sync.WaitGroup
	b.StartTimer()
	for index := 0; index < parallelism; index++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					panicch <- r
				}
			}()
			successes := 0
			i := index
			work := func() bool {
				if i < b.N {
					successes++
					i += parallelism
					return true
				}
				return false
			}
			err := worker(ctx, work)
			if err != nil {
				errch <- err
				return
			}
			atomic.AddUint64(&totalSuccesses, uint64(successes))
		}(index)
	}
	wg.Wait()
	b.StopTimer()
	select {
	case err := <-errch:
		cancel()
		b.Fatal(err)
	case r := <-panicch:
		cancel()
		panic(r)
	default:
		if totalSuccesses != uint64(b.N) {
			b.Fatalf("successes (%d) != b.N (%d)", totalSuccesses, b.N)
		}
	}
}
