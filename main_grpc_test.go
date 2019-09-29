package main

import (
	"crypto/tls"
	"log"
	"os"
	"sync/atomic"
	"testing"

	"golang.org/x/net/context"
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
	inParallel(parallelism, func(index int) {
		for i := index; i < b.N; i += parallelism {
			reply, err := client.SetInfo(context.Background(), &InfoRequest{
				Name:   "test",
				Age:    1,
				Height: 1,
			})
			if err != nil {
				b.Fatal(err)
			}
			if !reply.Success {
				b.Fatal("call failed")
			}
		}
	})
	b.StopTimer()
}

func BenchmarkGRPCSetInfoLoopback(b *testing.B) {
	loopbackAddr := "localhost:4443"
	s := mainGRPC(loopbackAddr, tlsCreds{certFile: "cert.pem", keyFile: "key.pem"})
	benchmarkGRPCSetInfo(b, loopbackAddr, 1)
	s.Stop()
}

func BenchmarkGRPCSetInfoLoopback16x(b *testing.B) {
	loopbackAddr := "localhost:4443"
	s := mainGRPC(loopbackAddr, tlsCreds{certFile: "cert.pem", keyFile: "key.pem"})
	benchmarkGRPCSetInfo(b, loopbackAddr, 16)
	s.Stop()
}

func BenchmarkGRPCSetInfoRemote(b *testing.B) {
	remoteAddr := os.Getenv("GRPC_REMOTE_ADDR")
	benchmarkGRPCSetInfo(b, remoteAddr, 1)
}

func BenchmarkGRPCSetInfoRemote16x(b *testing.B) {
	remoteAddr := os.Getenv("GRPC_REMOTE_ADDR")
	benchmarkGRPCSetInfo(b, remoteAddr, 16)
}

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
	ctx := context.Background()
	inParallel(parallelism, func(index int) {
		call, err := client.SetInfoStream(ctx)
		if err != nil {
			b.Fatal(err)
		}
		var t uint64 = 0
		for i := index; i < b.N; i += parallelism {
			if err := call.Send(&InfoRequest{
				Name:   "test",
				Age:    1,
				Height: 1,
			}); err != nil {
				b.Fatal(err)
			}
			reply, err := call.Recv()
			if err != nil {
				b.Fatal(err)
			}
			if !reply.Success {
				b.Fatal("call failed")
			}
			t++
		}
		atomic.AddUint64(&successes, t)
		if err := call.CloseSend(); err != nil {
			b.Fatal(err)
		}
	})
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
	ctx := context.Background()
	inParallel(parallelism, func(index int) {
		call, err := client.SetInfoStream(ctx)
		if err != nil {
			b.Fatal(err)
		}
		done := make(chan struct{})
		go func() {
			var t uint64 = 0
			for i := index; i < b.N; i += parallelism {
				reply, err := call.Recv()
				if err != nil {
					b.Fatal(err)
				}
				if !reply.Success {
					b.Fatal("call failed")
				}
				t++
			}
			atomic.AddUint64(&successes, t)
			close(done)
		}()
		for i := index; i < b.N; i += parallelism {
			if err := call.Send(&InfoRequest{
				Name:   "test",
				Age:    1,
				Height: 1,
			}); err != nil {
				b.Fatal(err)
			}
		}
		<-done
		if err := call.CloseSend(); err != nil {
			b.Fatal(err)
		}
	})
	if atomic.LoadUint64(&successes) != uint64(b.N) {
		b.Fatalf("successes (%d) != b.N (%d)", successes, b.N)
	}
	b.StopTimer()
}

func BenchmarkGRPCSetInfoAsyncStreamLoopback(b *testing.B) {
	loopbackAddr := "localhost:4443"
	s := mainGRPC(loopbackAddr, tlsCreds{certFile: "cert.pem", keyFile: "key.pem"})
	benchmarkGRPCSetInfoAsyncStream(b, loopbackAddr, 1)
	s.Stop()
}

func BenchmarkGRPCSetInfoAsyncStreamLoopback16x(b *testing.B) {
	loopbackAddr := "localhost:4443"
	s := mainGRPC(loopbackAddr, tlsCreds{certFile: "cert.pem", keyFile: "key.pem"})
	benchmarkGRPCSetInfoAsyncStream(b, loopbackAddr, 16)
	s.Stop()
}

func BenchmarkGRPCSetInfoAsyncStreamRemote(b *testing.B) {
	remoteAddr := os.Getenv("GRPC_REMOTE_ADDR")
	benchmarkGRPCSetInfoAsyncStream(b, remoteAddr, 1)
}

func BenchmarkGRPCSetInfoAsyncStreamRemote16x(b *testing.B) {
	remoteAddr := os.Getenv("GRPC_REMOTE_ADDR")
	benchmarkGRPCSetInfoAsyncStream(b, remoteAddr, 16)
}
