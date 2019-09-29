package main

import (
	"context"
	"os"
	"testing"

	"github.com/pkg/errors"
)

var benchmarkGRPCSetInfo = benchmarkGRPC(
	func(ctx context.Context, client InfoServerClient, work func() bool) error {
		for work() {
			reply, err := client.SetInfo(ctx, &InfoRequest{Name: "test", Age: 1, Height: 1})
			if err != nil {
				return errors.WithStack(err)
			}
			if !reply.Success {
				return errors.Errorf("call failed")
			}
		}
		return nil
	},
)

var benchmarkGRPCSetInfoLoopback = loopbackBenchmark(
	grpcPortBase, loopbackGRPC, benchmarkGRPCSetInfo)

func BenchmarkGRPCSetInfoLoopback(b *testing.B) {
	benchmarkGRPCSetInfoLoopback(b, 0, 1)
}

func BenchmarkGRPCSetInfoLoopback16x(b *testing.B) {
	benchmarkGRPCSetInfoLoopback(b, 1, 16)
}

func BenchmarkGRPCSetInfoRemote(b *testing.B) {
	benchmarkGRPCSetInfo(b, os.Getenv("GRPC_REMOTE_ADDR"), 1)
}

func BenchmarkGRPCSetInfoRemote16x(b *testing.B) {
	benchmarkGRPCSetInfo(b, os.Getenv("GRPC_REMOTE_ADDR"), 16)
}
