package main

import (
	"context"
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var benchmarkGRPCSetInfoStream = benchmarkGRPC(
	func(ctx context.Context, client InfoServerClient, work func() bool) error {
		call, err := client.SetInfoStream(ctx)
		if err != nil {
			return errors.WithStack(err)
		}
		for work() {
			if err := call.Send(&InfoRequest{Name: "test", Age: 1, Height: 1}); err != nil {
				return errors.WithStack(err)
			}
			reply, err := call.Recv()
			if err != nil {
				return errors.WithStack(err)
			}
			if !reply.Success {
				if err := call.CloseSend(); err != nil {
					logrus.Error(errors.WithStack(err))
				}
				return errors.New("call failed")
			}
		}
		return call.CloseSend()
	},
)

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
