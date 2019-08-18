package main

import (
	"crypto/tls"
	"errors"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"golang.org/x/net/context"
)

func mainGRPC(addr string, creds tlsCreds) *grpc.Server {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	config := &tls.Config{}
	cert, err := tls.LoadX509KeyPair(creds.certFile, creds.keyFile)
	if err != nil {
		log.Fatal(err)
	}
	config.Certificates = []tls.Certificate{cert}
	s := grpc.NewServer(grpc.Creds(credentials.NewTLS(config)))

	RegisterInfoServerServer(s, &server{})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatal(err)
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
