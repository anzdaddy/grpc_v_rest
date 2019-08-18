package main

import (
	"crypto/tls"
	"log"
	"os"
	"sync"
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
	var wg sync.WaitGroup
	for j := 0; j < parallelism; j++ {
		go func(j int) {
			for i := j; i < b.N; i += parallelism {
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
			wg.Done()
		}(j)
		wg.Add(1)
	}
	wg.Wait()
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

func benchmarkGRPCSetInfos(b *testing.B, addr string, parallelism int) {
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
	var wg sync.WaitGroup
	var successes uint64 = 0
	ctx := context.Background()
	for j := 0; j < parallelism; j++ {
		go func(j int) {
			call, err := client.SetInfos(ctx)
			if err != nil {
				b.Fatal(err)
			}
			done := make(chan struct{})
			go func() {
				var t uint64 = 0
				for i := j; i < b.N; i += parallelism {
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
			for i := j; i < b.N; i += parallelism {
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
			wg.Done()
		}(j)
		wg.Add(1)
	}
	wg.Wait()
	if atomic.LoadUint64(&successes) != uint64(b.N) {
		b.Fatalf("successes (%d) != b.N (%d)", successes, b.N)
	}
	b.StopTimer()
}

func BenchmarkGRPCSetInfosLoopback(b *testing.B) {
	loopbackAddr := "localhost:4443"
	s := mainGRPC(loopbackAddr, tlsCreds{certFile: "cert.pem", keyFile: "key.pem"})
	benchmarkGRPCSetInfos(b, loopbackAddr, 1)
	s.Stop()
}

func BenchmarkGRPCSetInfosLoopback16x(b *testing.B) {
	loopbackAddr := "localhost:4443"
	s := mainGRPC(loopbackAddr, tlsCreds{certFile: "cert.pem", keyFile: "key.pem"})
	benchmarkGRPCSetInfos(b, loopbackAddr, 16)
	s.Stop()
}

func BenchmarkGRPCSetInfosRemote(b *testing.B) {
	remoteAddr := os.Getenv("GRPC_REMOTE_ADDR")
	benchmarkGRPCSetInfos(b, remoteAddr, 1)
}

func BenchmarkGRPCSetInfosRemote16x(b *testing.B) {
	remoteAddr := os.Getenv("GRPC_REMOTE_ADDR")
	benchmarkGRPCSetInfos(b, remoteAddr, 16)
}
