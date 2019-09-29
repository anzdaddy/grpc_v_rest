package main

import (
	"flag"
	"fmt"

	"github.com/sirupsen/logrus"
)

type tlsCreds struct {
	certFile, keyFile string
}

func main() {
	grpcAddr := flag.String("grpc", ":4443", "gRPC listening address")
	restAddr := flag.String("rest", ":4444", "ReST listening address")
	var flagCreds tlsCreds
	flag.StringVar(&flagCreds.certFile, "cert", "cert.pem", "TLS Cert File")
	flag.StringVar(&flagCreds.keyFile, "key", "key.pem", "TLS Key File")
	flag.Parse()

	mainGRPC(*grpcAddr, flagCreds)
	logrus.Infof("gRPC %s", *grpcAddr)

	mainREST(*restAddr, flagCreds)
	logrus.Infof("REST %s", *restAddr)

	<-make(chan struct{})
}

// Validatable interface to describe what is validatable
type Validatable interface {
	Validate() error
}

// validationErrors - custom error type that can take many errors
type validationErrors []error

// Error - implementation of error interface
func (ve validationErrors) Error() string {
	var errStr string
	for _, v := range ve {
		errStr += fmt.Sprintf("%s\n", v.Error())
	}
	return errStr
}

func validate(v Validatable) error {
	return v.Validate()
}
