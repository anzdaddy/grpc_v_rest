package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"
)

func benchmarkRESTSetInfo(b *testing.B, addr string) {
	addr = "https://" + addr + "/info"
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	// run http posts against it
	b.StartTimer()
	for i := 0; i < b.N; i++ {
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
	b.StopTimer()
}

func BenchmarkRESTSetInfoLoopback(b *testing.B) {
	loopbackAddr := "localhost:4444"
	s := mainREST(loopbackAddr, tlsCreds{certFile: "cert.pem", keyFile: "key.pem"})
	benchmarkRESTSetInfo(b, loopbackAddr)
	s.Close()
}

func BenchmarkRESTSetInfoRemote(b *testing.B) {
	remoteAddr := os.Getenv("REST_REMOTE_ADDR")
	benchmarkRESTSetInfo(b, remoteAddr)
}
