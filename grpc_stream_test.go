package main

import (
	"crypto/tls"
	"log"
	"os"
	"sync/atomic"
	"testing"

	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func benchmarkGRPCSetInfoStream(b *testing.B, addr string, parallelism int) {
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
	if err := inParallel(context.Background(), parallelism, func(ctx context.Context, index int) error {
		call, err := client.SetInfoStream(ctx)
		if err != nil {
			return err
		}
		var t uint64 = 0
		for i := index; i < b.N; i += parallelism {
			if err := call.Send(&InfoRequest{
				Name:   "test",
				Age:    1,
				Height: 1,
			}); err != nil {
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

func BenchmarkGRPCSetInfoStreamLoopback(b *testing.B) {
	loopbackAddr := "localhost:4443"
	s := mainGRPC(loopbackAddr, tlsCreds{certFile: "cert.pem", keyFile: "key.pem"})
	benchmarkGRPCSetInfoStream(b, loopbackAddr, 1)
	s.Stop()
}

func BenchmarkGRPCSetInfoStreamLoopback16x(b *testing.B) {
	loopbackAddr := "localhost:4443"
	s := mainGRPC(loopbackAddr, tlsCreds{certFile: "cert.pem", keyFile: "key.pem"})
	benchmarkGRPCSetInfoStream(b, loopbackAddr, 16)
	s.Stop()
}

func BenchmarkGRPCSetInfoStreamRemote(b *testing.B) {
	remoteAddr := os.Getenv("GRPC_REMOTE_ADDR")
	benchmarkGRPCSetInfoStream(b, remoteAddr, 1)
}

func BenchmarkGRPCSetInfoStreamRemote16x(b *testing.B) {
	remoteAddr := os.Getenv("GRPC_REMOTE_ADDR")
	benchmarkGRPCSetInfoStream(b, remoteAddr, 16)
}
