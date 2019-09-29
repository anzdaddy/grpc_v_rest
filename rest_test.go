package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func benchmarkRESTSetInfo(b *testing.B, addr string, parallelism int) {
	url := fmt.Sprintf("https://%s/info", addr)
	parallelBenchmark(b, parallelism, func(ctx context.Context, work <-chan int) error {
		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}
		for range work {
			reqData, err := json.Marshal(apiInput{Name: "test", Age: 1, Height: 1})
			if err != nil {
				return errors.WithStack(err)
			}
			req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqData))
			if err != nil {
				return errors.WithStack(err)
			}
			req.Header.Set("Content-Type", "application/json")
			resp, err := client.Do(req)
			if err != nil {
				return errors.WithStack(err)
			}
			var r apiResponse
			if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
				if err2 := resp.Body.Close(); err2 != nil {
					logrus.Error(errors.WithStack(err2))
				}
				return errors.Wrap(err, "Error parsing JSON")
			}
			if !r.Success {
				return errors.New("call failed")
			}
		}
		return nil
	})
}

var benchmarkRESTSetInfoLoopback = loopbackBenchmark(
	restPortBase, loopbackREST, benchmarkRESTSetInfo)

func BenchmarkRESTSetInfoLoopback(b *testing.B) {
	benchmarkRESTSetInfoLoopback(b, 0, 1)
}

func BenchmarkRESTSetInfoLoopback16x(b *testing.B) {
	benchmarkRESTSetInfoLoopback(b, 1, 16)
}

func BenchmarkRESTSetInfoRemote(b *testing.B) {
	benchmarkRESTSetInfo(b, os.Getenv("REST_REMOTE_ADDR"), 1)
}

func BenchmarkRESTSetInfoRemote16x(b *testing.B) {
	benchmarkRESTSetInfo(b, os.Getenv("REST_REMOTE_ADDR"), 16)
}
