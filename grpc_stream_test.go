package main

import (
	"context"
	"os"
	"testing"

	"github.com/pkg/errors"
)

var benchmarkGRPCSetInfoStream = benchmarkGRPC(
	func(ctx context.Context, client InfoServerClient, work <-chan int) error {
		call, err := client.SetInfoStream(ctx)
		if err != nil {
			return err
		}
		for range work {
			if err := call.Send(&InfoRequest{Name: "test", Age: 1, Height: 1}); err != nil {
				return err
			}
			reply, err := call.Recv()
			if err != nil {
				return err
			}
			if !reply.Success {
				call.CloseSend()
				return errors.Errorf("call failed")
			}
		}
		return call.CloseSend()
	},
)

var benchmarkGRPCSetInfoStreamLoopback = loopbackBenchmark(
	grpcStreamPortBase, loopbackGRPC, benchmarkGRPCSetInfoStream)

func BenchmarkGRPCSetInfoStreamLoopback(b *testing.B) {
	benchmarkGRPCSetInfoStreamLoopback(b, 0, 1)
}

func BenchmarkGRPCSetInfoStreamLoopback16x(b *testing.B) {
	benchmarkGRPCSetInfoStreamLoopback(b, 1, 16)
}

func BenchmarkGRPCSetInfoStreamRemote(b *testing.B) {
	benchmarkGRPCSetInfoStream(b, os.Getenv("GRPC_REMOTE_ADDR"), 1)
}

func BenchmarkGRPCSetInfoStreamRemote16x(b *testing.B) {
	benchmarkGRPCSetInfoStream(b, os.Getenv("GRPC_REMOTE_ADDR"), 16)
}
