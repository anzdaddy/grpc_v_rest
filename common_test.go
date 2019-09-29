package main

import "fmt"

const (
	grpcUnaryPortBase = 4440 + 10*iota
	grpcStreamPortBase
	grpcAsyncStreamPortBase
	restPortBase
)

func loopbackTestAddress(port int) string {
	return fmt.Sprintf("localhost:%d", port)
}
