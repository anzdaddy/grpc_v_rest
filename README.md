# GRPC versus REST+HTTP/2

### Loopback test:

```bash
go test -bench Loopback
```

### Local test

Shell 1:

```
go build
./makecerts.sh
go run . -cert cert.pem -key key.pem
```

Shell 2:

```
GRPC_REMOTE_ADDR=localhost:4443 REST_REMOTE_ADDR=localhost:4444 go test -bench Remote
```

### Remote test

Machine 1:

```
go build
./makecerts.sh
go run . -cert cert.pem -key key.pem
```

Machine 2:

```
SRV=<<machine1-dns-or-ip>>
GRPC_REMOTE_ADDR=${SRV}:4443 REST_REMOTE_ADDR=${SRV}:4444 go test -bench Remote -benchtime 10s
```

## Results

The following results track four techniques for implementing a request/response flow.

1. **gRPC** uses a single-in/single-out gRPC call per request.
2. **gRPC streaming** uses a single stream-in/stream-out gRPC call for every request.
3. **gRPC async streaming** use a single stream-in/stream-out gRPC call, as above. Requests are sent over one goroutine without waiting for responses. Responses are harvested independently by a separate goroutine.
4. **REST** calls a RESTful API over HTTP/2 once per request.

For each technique, we run the client on a single-goroutine (1x) and on 16 goroutines (16x).

### Loopback (in-process server)

```
$ go test -bench Loopback
goos: darwin
goarch: amd64
pkg: github.com/anzdaddy/grpc_v_rest
BenchmarkGRPCSetInfoLoopback-8                 	   10000	    193134 ns/op
BenchmarkGRPCSetInfoLoopback16x-8              	   50000	     24905 ns/op
BenchmarkGRPCSetInfoStreamLoopback-8           	ERRO[0004] rpc error: code = Canceled desc = context canceled
   10000	    120618 ns/op
BenchmarkGRPCSetInfoStreamLoopback16x-8        	ERRO[0006] rpc error: code = Canceled desc = context canceled
  100000	     12562 ns/op
BenchmarkGRPCSetInfoAsyncStreamLoopback-8      	ERRO[0007] rpc error: code = Canceled desc = context canceled
ERRO[0007] rpc error: code = Canceled desc = context canceled
  200000	      5168 ns/op
BenchmarkGRPCSetInfoAsyncStreamLoopback16x-8   	ERRO[0008] rpc error: code = Canceled desc = context canceled
ERRO[0008] rpc error: code = Canceled desc = context canceled
ERRO[0008] rpc error: code = Canceled desc = context canceled
  500000	      2920 ns/op
ERRO[0010] rpc error: code = Canceled desc = context canceled
BenchmarkRESTSetInfoLoopback-8                 	   10000	    176539 ns/op
BenchmarkRESTSetInfoLoopback16x-8              	   20000	     54953 ns/op
PASS
ok  	github.com/anzdaddy/grpc_v_rest	14.649s
```

¹Something's NQR with StreamLoopback scenarios. Further investigation is needed.

![](https://chart.googleapis.com/chart?cht=bvg&chs=300x300&chdl=gRPC|REST&chd=t:193.134,24.905|176.539,54.953&chds=a&chxt=x,y&chxl=0:|1x|16x&chco=A03333,4D89F9&chxs=1N**+µs&chma=10,10,10,10&chbh=30,5,20 "gRPC vs REST loopback")

### Local (“Remote” to localhost)

```
$ GRPC_REMOTE_ADDR=localhost:4443 REST_REMOTE_ADDR=localhost:4444 go test -bench Remote
goos: darwin
goarch: amd64
pkg: github.com/anzdaddy/grpc_v_rest
BenchmarkGRPCSetInfoRemote-8                 	    5000	    271194 ns/op
BenchmarkGRPCSetInfoRemote16x-8              	   50000	     31229 ns/op
BenchmarkGRPCSetInfoStreamRemote-8           	   10000	    182109 ns/op
BenchmarkGRPCSetInfoStreamRemote16x-8        	  100000	     16459 ns/op
BenchmarkGRPCSetInfoAsyncStreamRemote-8      	  200000	      5543 ns/op
BenchmarkGRPCSetInfoAsyncStreamRemote16x-8   	  500000	      3643 ns/op
BenchmarkRESTSetInfoRemote-8                 	   10000	    235802 ns/op
BenchmarkRESTSetInfoRemote16x-8              	   20000	     68973 ns/op
PASS
ok  	github.com/anzdaddy/grpc_v_rest	17.719s
```

²They're basically neck and neck. Different runs yield different winners.

![](https://chart.googleapis.com/chart?cht=bvg&chs=500x300&chdl=gRPC|gRPC+(stream)|gRPC+(async+stream)|REST&chd=t:271.194,31.229|182.109,16.459|5.543,3.643|235.802,68.973&chds=a&chxt=x,y&chxl=0:|1x|16x&chco=A03333,C09999,FF6600,4D89F9&chxs=1N**+µs&chma=10,10,10,10&chbh=30,5,20 "gRPC vs REST loopback")

### Remote (One physical machine to another on a 1GB LAN)

```
$ GRPC_REMOTE_ADDR=192.168.87.42:4443 REST_REMOTE_ADDR=192.168.87.42:4444 go test -bench Remote
goos: darwin
goarch: amd64
pkg: github.com/anzdaddy/grpc_v_rest
BenchmarkGRPCSetInfoRemote-4                          300       3849929 ns/op
BenchmarkGRPCSetInfoRemote16x-4                      5000        270626 ns/op
BenchmarkGRPCSetInfoStreamRemote-4                    500       3599877 ns/op
BenchmarkGRPCSetInfoStreamRemote16x-4               10000        244286 ns/op
BenchmarkGRPCSetInfoAsyncStreamRemote-4            500000          5234 ns/op
BenchmarkGRPCSetInfoAsyncStreamRemote16x-4         500000          3019 ns/op
BenchmarkRESTSetInfoRemote-4                          300       3580039 ns/op
BenchmarkRESTSetInfoRemote16x-4                      3000        346382 ns/op
PASS
ok      github.com/anzdaddy/grpc_v_rest    16.470s
```

³Hmm, async streaming is even faster over the wire!

![](https://chart.googleapis.com/chart?cht=bvg&chs=500x300&chdl=gRPC|gRPC+(stream)|gRPC+(async+stream)|REST&chd=t:3.849929,.270626|3.599877,.244286|.005234,.003019|3.580039,.346382&chds=a&chxt=x,y&chxl=0:|1x|16x&chco=A03333,C09999,FF6600,4D89F9&chma=10,10,10,10&chbh=30,5,20&chxs=1N**+ms "gRPC vs REST loopback")

## Selected observations

1. gRPC scales much better than REST under 16x parallelised clients, clocking in at more than 2x faster.
2. On localhost, streaming gRPC request/response flows are much faster than non-streaming gRPC calls, but the difference largely goes away over a LAN.
3. Async streaming offers vastly more throughput than any other technique. At 16x parallelism, it is 80x faster than the next best option.

## Machine specs

### Machine 1 (also used for loopback and local):

```
$ system_profiler -detailLevel mini SPHardwareDataType
Hardware:

    Hardware Overview:

      Model Name: MacBook Pro
      Model Identifier: MacBookPro10,1
      Processor Name: Intel Core i7
      Processor Speed: 2.6 GHz
      Number of Processors: 1
      Total Number of Cores: 4
      L2 Cache (per Core): 256 KB
      L3 Cache: 6 MB
      Hyper-Threading Technology: Enabled
      Memory: 16 GB
      Boot ROM Version: 257.0.0.0.0
      SMC Version (system): 2.3f36
```

### Machine 2:

```
$ system_profiler -detailLevel mini SPHardwareDataType
Hardware:

    Hardware Overview:

      Model Name: MacBook Pro
      Model Identifier: MacBookPro13,2
      Processor Name: Intel Core i5
      Processor Speed: 3.1 GHz
      Number of Processors: 1
      Total Number of Cores: 2
      L2 Cache (per Core): 256 KB
      L3 Cache: 4 MB
      Memory: 16 GB
      Boot ROM Version: 254.0.0.0.0
      SMC Version (system): 2.37f20
```
