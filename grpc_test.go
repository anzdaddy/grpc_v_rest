package main

import (
	"context"
	"crypto/tls"
	"log"
	"os"
	"testing"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func benchmarkGRPCSetInfo(b *testing.B, addr string, parallelism int) {
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
	if err := inParallel(context.Background(), parallelism, func(ctx context.Context, index int) error {
		for i := index; i < b.N; i += parallelism {
			reply, err := client.SetInfo(ctx, &InfoRequest{
				Name:   "test",
				Age:    1,
				Height: 1,
			})
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

func BenchmarkGRPCSetInfoLoopback(b *testing.B) {
	loopbackAddr := loopbackTestAddress(grpcUnaryPortBase + 0)
	defer mainGRPC(loopbackAddr, tlsCreds{certFile: "cert.pem", keyFile: "key.pem"}).Stop()
	time.Sleep(time.Millisecond)
	benchmarkGRPCSetInfo(b, loopbackAddr, 1)
}

func BenchmarkGRPCSetInfoLoopback16x(b *testing.B) {
	loopbackAddr := loopbackTestAddress(grpcUnaryPortBase + 1)
	defer mainGRPC(loopbackAddr, tlsCreds{certFile: "cert.pem", keyFile: "key.pem"}).Stop()
	time.Sleep(time.Millisecond)
	benchmarkGRPCSetInfo(b, loopbackAddr, 16)
}

func BenchmarkGRPCSetInfoRemote(b *testing.B) {
	remoteAddr := os.Getenv("GRPC_REMOTE_ADDR")
	benchmarkGRPCSetInfo(b, remoteAddr, 1)
}

func BenchmarkGRPCSetInfoRemote16x(b *testing.B) {
	remoteAddr := os.Getenv("GRPC_REMOTE_ADDR")
	benchmarkGRPCSetInfo(b, remoteAddr, 16)
}