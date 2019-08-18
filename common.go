package main

import (
	"flag"
	"fmt"

	"github.com/sirupsen/logrus"
)

type tlsCreds struct {
	certFile, keyFile string
}

var (
	flagCreds tlsCreds
)

func init() {
	flag.StringVar(&flagCreds.certFile, "cert", "test.crt", "TLS Cert File")
	flag.StringVar(&flagCreds.keyFile, "key", "test.key", "TLS Key File")
	flag.Parse()
}

func main() {
	grpcAddr := ":4443"
	mainGRPC(grpcAddr, flagCreds)
	logrus.Infof("gRPC %s", grpcAddr)

	restAddr := ":4444"
	mainREST(restAddr, flagCreds)
	logrus.Infof("REST %s", restAddr)

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
