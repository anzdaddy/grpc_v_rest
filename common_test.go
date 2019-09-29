package main

import (
	"crypto/tls"
	"fmt"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	grpcPortBase = 4440 + 10*iota
	grpcStreamPortBase
	grpcAsyncStreamPortBase
	restPortBase
)

func loopbackTestAddress(port int) string {
	return fmt.Sprintf("localhost:%d", port)
}

func loopbackTestCreds() tlsCreds {
	return tlsCreds{
		certFile: "cert.pem",
		keyFile:  "key.pem",
	}
}

func loopbackBenchmark(
	portBase int,
	startServer func(addr string) (cancel func()),
	benchmark func(b *testing.B, addr string, parallelism int),
) func(b *testing.B, portOffset, parallelism int) {
	return func(b *testing.B, portOffset, parallelism int) {
		loopbackAddr := loopbackTestAddress(portBase + portOffset)
		defer startServer(loopbackAddr)
		time.Sleep(10 * time.Millisecond)
		benchmark(b, loopbackAddr, parallelism)
	}
}

func loopbackGRPC(loopbackAddr string) (cancel func()) {
	return mainGRPC(loopbackAddr, loopbackTestCreds()).Stop
}

func loopbackReST(loopbackAddr string) (cancel func()) {
	server := mainREST(loopbackAddr, loopbackTestCreds())
	return func() { server.Close() }
}

func grpcSetInfoClient(addr string) (conn *grpc.ClientConn, client InfoServerClient, err error) {
	config := &tls.Config{}
	config.InsecureSkipVerify = true
	// Set up a connection to the server.
	conn, err = grpc.Dial(addr, grpc.WithTransportCredentials(credentials.NewTLS(config)))
	if err != nil {
		return nil, nil, err
	}
	return conn, NewInfoServerClient(conn), nil
}
