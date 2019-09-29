package main

import (
	"context"
	"os"
	"testing"

	"github.com/pkg/errors"
)

var benchmarkGRPCSetInfoAsyncStream = benchmarkGRPC(
	func(ctx context.Context, client InfoServerClient, work func() bool) error {
		call, err := client.SetInfoStream(ctx)
		if err != nil {
			return err
		}
		done := make(chan error, 1)
		defer func() {
			if err2 := <-done; err2 != nil {
				err = err2
			}
		}()
		sent := make(chan struct{}, 1000)
		go func() {
			defer close(done)
			for range sent {
				reply, err := call.Recv()
				if err != nil {
					done <- err
					return
				}
				if !reply.Success {
					done <- errors.Errorf("call failed")
					return
				}
			}
		}()
		for work() {
			if err := call.Send(&InfoRequest{Name: "test", Age: 1, Height: 1}); err != nil {
				return err
			}
			sent <- struct{}{}
		}
		close(sent)
		if err := call.CloseSend(); err != nil {
			return err
		}
		return nil
	},
)

var benchmarkGRPCSetInfoAsyncStreamLoopback = loopbackBenchmark(
	grpcAsyncStreamPortBase, loopbackGRPC, benchmarkGRPCSetInfoAsyncStream)

func BenchmarkGRPCSetInfoAsyncStreamLoopback(b *testing.B) {
	benchmarkGRPCSetInfoAsyncStreamLoopback(b, 0, 1)
}

func BenchmarkGRPCSetInfoAsyncStreamLoopback16x(b *testing.B) {
	benchmarkGRPCSetInfoAsyncStreamLoopback(b, 1, 16)
}

func BenchmarkGRPCSetInfoAsyncStreamRemote(b *testing.B) {
	benchmarkGRPCSetInfoAsyncStream(b, os.Getenv("GRPC_REMOTE_ADDR"), 1)
}

func BenchmarkGRPCSetInfoAsyncStreamRemote16x(b *testing.B) {
	benchmarkGRPCSetInfoAsyncStream(b, os.Getenv("GRPC_REMOTE_ADDR"), 16)
}
