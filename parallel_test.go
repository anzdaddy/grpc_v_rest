package main

import (
	"context"
	"sync"
)

func inParallel(ctx context.Context, parallelism int, work func(ctx context.Context, index int) error) error {
	errch := make(chan error, parallelism)
	panicch := make(chan interface{}, parallelism)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	var wg sync.WaitGroup
	for index := 0; index < parallelism; index++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					panicch <- r
				}
			}()
			if err := work(ctx, index); err != nil {
				errch <- err
				return
			}
		}(index)
	}
	wg.Wait()
	select {
	case err := <-errch:
		cancel()
		return err
	case r := <-panicch:
		cancel()
		panic(r)
	default:
		return nil
	}
}
