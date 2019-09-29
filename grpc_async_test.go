package main

import (
	"context"
	"crypto/tls"
	"log"
	"os"
	"sync/atomic"
	"testing"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func benchmarkGRPCSetInfoAsyncStream(b *testing.B, addr string, parallelism int) {
	config := &tls.Config{}
	config.InsecureSkipVerify = true
	// Set up a connection to the server.
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(credentials.NewTLS(config)))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := NewInfoServerClient(conn)

	// run grpc calls against it
	b.StartTimer()
	var successes uint64 = 0
	if err := inParallel(context.Background(), parallelism, func(ctx context.Context, index int) (err error) {
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
		go func() {
			defer close(done)
			var t uint64 = 0
			for i := index; i < b.N; i += parallelism {
				reply, err := call.Recv()
				if err != nil {
					done <- err
					return
				}
				if !reply.Success {
					done <- errors.Errorf("call failed")
					return
				}
				t++
			}
			atomic.AddUint64(&successes, t)
		}()
		for i := index; i < b.N; i += parallelism {
			if err := call.Send(&InfoRequest{
				Name:   "test",
				Age:    1,
				Height: 1,
			}); err != nil {
				return err
			}
		}
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
