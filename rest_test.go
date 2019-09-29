package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/pkg/errors"
)

func benchmarkRESTSetInfo(b *testing.B, addr string, parallelism int) {
	url := "https://" + addr + "/info"
	b.StartTimer()
	if err := inParallel(context.Background(), parallelism, func(ctx context.Context, j int) error {
		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}
		for i := j; i < b.N; i += parallelism {
			reqData, err := json.Marshal(apiInput{
				Name:   "test",
				Age:    1,
				Height: 1,
			})
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
	}); err != nil {
		b.Fatal(err)
	}
	b.StopTimer()
}

func BenchmarkRESTSetInfoLoopback(b *testing.B) {
	loopbackAddr := loopbackTestAddress(restPortBase + 0)
	defer mainREST(loopbackAddr, tlsCreds{certFile: "cert.pem", keyFile: "key.pem"}).Close()
	time.Sleep(time.Millisecond)
	benchmarkRESTSetInfo(b, loopbackAddr, 1)
}

func BenchmarkRESTSetInfoLoopback16x(b *testing.B) {
	loopbackAddr := loopbackTestAddress(restPortBase + 1)
	defer mainREST(loopbackAddr, tlsCreds{certFile: "cert.pem", keyFile: "key.pem"}).Close()
	time.Sleep(time.Millisecond)
	benchmarkRESTSetInfo(b, loopbackAddr, 16)
}

func BenchmarkRESTSetInfoRemote(b *testing.B) {
	remoteAddr := os.Getenv("REST_REMOTE_ADDR")
	benchmarkRESTSetInfo(b, remoteAddr, 1)
}

func BenchmarkRESTSetInfoRemote16x(b *testing.B) {
	remoteAddr := os.Getenv("REST_REMOTE_ADDR")
	benchmarkRESTSetInfo(b, remoteAddr, 16)
}
