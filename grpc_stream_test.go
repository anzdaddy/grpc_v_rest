package main

import (
	"context"
	"os"
	"sync/atomic"
	"testing"

	"github.com/pkg/errors"
)

func benchmarkGRPCSetInfoStream(b *testing.B, addr string, parallelism int) {
	conn, client, err := grpcSetInfoClient(addr)
	if err != nil {
		b.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	// run grpc calls against it
	b.StartTimer()
	var successes uint64 = 0
	if err := inParallel(context.Background(), parallelism, func(ctx context.Context, index int) error {
		call, err := client.SetInfoStream(ctx)
		if err != nil {
			return err
		}
		var t uint64 = 0
		for i := index; i < b.N; i += parallelism {
			if err := call.Send(&InfoRequest{Name: "test", Age: 1, Height: 1}); err != nil {
				return err
			}
			reply, err := call.Recv()
			if err != nil {
				return err
			}
			if !reply.Success {
				return errors.Errorf("call failed")
			}
			t++
		}
		atomic.AddUint64(&successes, t)
		if err := call.CloseSend(); err != nil {
			return err
		}
		return nil
	}); err != nil {
		b.Fatal(err)
	}
	if atomic.LoadUint64(&successes) != uint64(b.N) {
		b.Fatalf("successes (%d) != b.N (%d)", successes, b.N)
	}
	b.StopTimer()
}

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
