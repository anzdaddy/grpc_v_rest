package main

import (
	"context"
	"os"
	"testing"

	"github.com/pkg/errors"
)

func benchmarkGRPCSetInfo(b *testing.B, addr string, parallelism int) {
	conn, client, err := grpcSetInfoClient(addr)
	if err != nil {
		b.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	// run grpc calls against it
	b.StartTimer()
	if err := inParallel(context.Background(), parallelism, func(ctx context.Context, index int) error {
		for i := index; i < b.N; i += parallelism {
			reply, err := client.SetInfo(ctx, &InfoRequest{Name: "test", Age: 1, Height: 1})
			if err != nil {
				return errors.WithStack(err)
			}
			if !reply.Success {
				return errors.Errorf("call failed")
			}
		}
		return nil
	}); err != nil {
		b.Fatal(err)
	}
	b.StopTimer()
}

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
