package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	grpcPortBase = 4440 + 10*iota
	grpcStreamPortBase
	grpcAsyncStreamPortBase
	restPortBase
)

func loopbackTestCreds() tlsCreds {
	return tlsCreds{
		certFile: "cert.pem",
		keyFile:  "key.pem",
	}
}

func loopbackBenchmark(
	portBase int,
	startServer func(addr string) CloserFunc,
	benchmark func(b *testing.B, addr string, parallelism int),
) func(b *testing.B, portOffset, parallelism int) {
	return func(b *testing.B, portOffset, parallelism int) {
		loopbackAddr := fmt.Sprintf("localhost:%d", portBase+portOffset)
		logrus.SetReportCaller(true)
		defer logrus.SetReportCaller(false)
		defer startServer(loopbackAddr)()
		time.Sleep(10 * time.Millisecond)
		benchmark(b, loopbackAddr, parallelism)
	}
}

type CloserFunc func()

func loopbackGRPC(loopbackAddr string) CloserFunc {
	return mainGRPC(loopbackAddr, loopbackTestCreds()).Stop
}

func loopbackREST(loopbackAddr string) CloserFunc {
	server := mainREST(loopbackAddr, loopbackTestCreds())
	return func() {
		if err := server.Close(); err != nil {
			logrus.Error(errors.WithStack(err))
		}
	}
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

func benchmarkGRPC(
	worker func(ctx context.Context, client InfoServerClient, work <-chan int) error,
) func(b *testing.B, addr string, parallelism int) {
	return func(b *testing.B, addr string, parallelism int) {
		conn, client, err := grpcSetInfoClient(addr)
		if err != nil {
			b.Fatalf("failed to connect: %v", err)
		}
		defer conn.Close()
		parallelBenchmark(b, parallelism, func(ctx context.Context, work <-chan int) error {
			return worker(ctx, client, work)
		})
	}
}
