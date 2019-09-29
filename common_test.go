package main

import (
	"fmt"
	"time"
)

const (
	grpcUnaryPortBase = 4440 + 10*iota
	grpcStreamPortBase
	grpcAsyncStreamPortBase
	restPortBase
)

func loopbackTestAddress(port int) string {
	return fmt.Sprintf("localhost:%d", port)
}

func giveLoopbackServerTimeToStart() {
	time.Sleep(10 * time.Millisecond)
}
