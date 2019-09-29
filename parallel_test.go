package main

import (
	"sync"
)

func inParallel(parallelism int, work func(index int)) {
	var wg sync.WaitGroup
	for index := 0; index < parallelism; index++ {
		wg.Add(1)
		go func(index int) {
			work(index)
			wg.Done()
		}(index)
	}
	wg.Wait()
}
