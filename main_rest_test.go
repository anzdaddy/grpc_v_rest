package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"sync"
	"testing"
)

func benchmarkRESTSetInfo(b *testing.B, addr string, parallelism int) {
	addr = "https://" + addr + "/info"
	b.StartTimer()
	var wg sync.WaitGroup
	for j := 0; j < parallelism; j++ {
		go func(j int) {
			client := &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: true,
					},
				},
			}
			for i := j; i < b.N; i += parallelism {
				req, err := json.Marshal(apiInput{
					Name:   "test",
					Age:    1,
					Height: 1,
				})
				if err != nil {
					b.Fatal(err)
				}
				resp, err := client.Post(addr, "application/json", bytes.NewBuffer(req))
				if err != nil {
					b.Fatal(err)
				}
				defer resp.Body.Close()
				var r apiResponse
				var buf bytes.Buffer
				io.Copy(&buf, resp.Body)
				data := buf.String()
				if err := json.NewDecoder(&buf).Decode(&r); err != nil {
					b.Fatalf("Error parsing JSON: %s\n<<%s>>", err, data)
				}
				if !r.Success {
					b.Fatalf("call failed\n<<%s>>", data)
				}
			}
			wg.Done()
		}(j)
		wg.Add(1)
	}
	wg.Wait()
	b.StopTimer()
}

func BenchmarkRESTSetInfoLoopback(b *testing.B) {
	loopbackAddr := "localhost:4444"
	s := mainREST(loopbackAddr, tlsCreds{certFile: "cert.pem", keyFile: "key.pem"})
	benchmarkRESTSetInfo(b, loopbackAddr, 1)
	s.Close()
}

func BenchmarkRESTSetInfoLoopback16x(b *testing.B) {
	loopbackAddr := "localhost:4444"
	s := mainREST(loopbackAddr, tlsCreds{certFile: "cert.pem", keyFile: "key.pem"})
	benchmarkRESTSetInfo(b, loopbackAddr, 16)
	s.Close()
}

func BenchmarkRESTSetInfoRemote(b *testing.B) {
	remoteAddr := os.Getenv("REST_REMOTE_ADDR")
	benchmarkRESTSetInfo(b, remoteAddr, 1)
}

func BenchmarkRESTSetInfoRemote16x(b *testing.B) {
	remoteAddr := os.Getenv("REST_REMOTE_ADDR")
	benchmarkRESTSetInfo(b, remoteAddr, 16)
}
