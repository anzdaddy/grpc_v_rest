package main

import (
	"crypto/tls"
	"log"
	"os"
	"testing"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func benchmarkGRPCSetInfo(b *testing.B, addr string) {
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
	for i := 0; i < b.N; i++ {
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
	b.StopTimer()
}

func BenchmarkGRPCSetInfoLoopback(b *testing.B) {
	loopbackAddr := "localhost:4443"
	s := mainGRPC(loopbackAddr, tlsCreds{certFile: "cert.pem", keyFile: "key.pem"})
	benchmarkGRPCSetInfo(b, loopbackAddr)
	s.Stop()
}

func BenchmarkGRPCSetInfoRemote(b *testing.B) {
	remoteAddr := os.Getenv("GRPC_REMOTE_ADDR")
	benchmarkGRPCSetInfo(b, remoteAddr)
}
