package main

import (
	"context"
	"crypto/tls"
	"errors"
	"io"
	"net"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

//go:generate protoc --go_out=plugins=grpc:. info.proto

func mainGRPC(addr string, creds tlsCreds) *grpc.Server {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logrus.Fatalf("failed to listen: %v", err)
	}

	config := &tls.Config{}
	cert, err := tls.LoadX509KeyPair(creds.certFile, creds.keyFile)
	if err != nil {
		logrus.Fatal(err)
	}
	config.Certificates = []tls.Certificate{cert}
	s := grpc.NewServer(grpc.Creds(credentials.NewTLS(config)))

	RegisterInfoServerServer(s, &server{})
	go func() {
		if err := s.Serve(lis); err != nil {
			logrus.Fatal(err)
		}
	}()
	return s
}

type server struct{}

// SetInfo - implements our InfoServer
func (s *server) SetInfo(ctx context.Context, in *InfoRequest) (*InfoReply, error) {
	if err := validate(in); err != nil {
		return &InfoReply{
			Success: false,
			Reason:  err.Error(),
		}, err
	}
	return &InfoReply{
		Success: true,
	}, nil
}

// SetInfos implements the streaming model
func (s *server) SetInfoStream(server InfoServer_SetInfoStreamServer) error {
	for {
		in, err := server.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			logrus.Error(err)
			return err
		}
		if err := validate(in); err != nil {
			if err := server.Send(&InfoReply{
				Success: false,
				Reason:  err.Error(),
			}); err != nil {
				logrus.Error(err)
				return err
			}
			continue
		}
		if err := server.Send(&InfoReply{
			Success: true,
		}); err != nil {
			logrus.Error(err)
			return err
		}
	}
}

// Validate - implement validatable
func (ir *InfoRequest) Validate() error {
	var err validationErrors
	if ir.Name == "" {
		err = append(err, errors.New("Name must be present"))
	}
	if ir.Age <= 0 {
		err = append(err, errors.New("Age must be real"))
	}
	if ir.Height <= 0 {
		err = append(err, errors.New("Height must be real"))
	}
	if len(err) == 0 {
		return nil
	}
	return err
}
