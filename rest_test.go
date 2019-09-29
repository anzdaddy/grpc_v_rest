package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/pkg/errors"
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
				return err
			}
			req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqData))
			if err != nil {
				return err
			}
			req.Header.Set("Content-Type", "application/json")
			resp, err := client.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			var r apiResponse
			var buf bytes.Buffer
			if _, err := io.Copy(&buf, resp.Body); err != nil {
				return errors.Errorf("Error copying resp.Body: %s", err)
			}
			respData := buf.String()
			if err := json.NewDecoder(&buf).Decode(&r); err != nil {
				return errors.Errorf("Error parsing JSON: %s\n%s", err, respData)
			}
			if !r.Success {
				return errors.Errorf("call failed\n%s", respData)
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
