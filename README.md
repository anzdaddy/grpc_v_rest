# GRPC versus REST

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

The following results track three techniques — gRPC, gRPC streaming and REST. For
each, we run the client on a single-goroutine (1x) and on 16 goroutines (16x).

### Loopback (in-process server)

```
$ go test -v -bench Loopback
goos: darwin
goarch: amd64
pkg: github.com/anzdaddy/grpc_v_rest
BenchmarkGRPCSetInfoLoopback-8            	   10000	    182394 ns/op
BenchmarkGRPCSetInfoLoopback16x-8         	   50000	     25942 ns/op
BenchmarkGRPCSetInfoStreamLoopback-8      	ERRO[0004] rpc error: code = Canceled desc = context canceled
  200000	      5118 ns/op¹
BenchmarkGRPCSetInfoStreamLoopback16x-8   	ERRO[0006] rpc error: code = Canceled desc = context canceled
ERRO[0006] rpc error: code = Canceled desc = context canceled
  500000	      2811 ns/op¹
BenchmarkRESTSetInfoLoopback-8            	   10000	    176904 ns/op
BenchmarkRESTSetInfoLoopback16x-8         	   20000	     55124 ns/op
PASS
ok  	github.com/anzdaddy/grpc_v_rest	11.757s
```

¹Something's obviously NQR with in-process streaming calls. Further investigation is needed.

![](https://chart.googleapis.com/chart?cht=bvg&chs=300x300&chdl=gRPC|REST&chd=t:182.394,25.942|176.904,55.124&chds=a&chxt=x,y&chxl=0:|1x|16x|&chco=4D89F9,C6D9FD&chxs=1N**+µs&chma=10,10,10,10&chbh=30,5,20 "gRPC vs REST loopback")

### Local (“Remote” to localhost)

```
$ GRPC_REMOTE_ADDR=localhost:4443 REST_REMOTE_ADDR=localhost:4444 go test -bench Remote
goos: darwin
goarch: amd64
pkg: github.com/anzdaddy/grpc_v_rest
BenchmarkGRPCSetInfoRemote-8            	    5000	    259356 ns/op²
BenchmarkGRPCSetInfoRemote16x-8         	   50000	     32297 ns/op
BenchmarkGRPCSetInfoStreamRemote-8      	  200000	      5117 ns/op
BenchmarkGRPCSetInfoStreamRemote16x-8   	  500000	      3448 ns/op
BenchmarkRESTSetInfoRemote-8            	   10000	    260499 ns/op²
BenchmarkRESTSetInfoRemote16x-8         	   20000	     58950 ns/op
PASS
ok  	github.com/anzdaddy/grpc_v_rest	12.424s
```

²They're basically neck and neck. Different runs yield different winners.

![](https://chart.googleapis.com/chart?cht=bvg&chs=450x300&chdl=gRPC|gRPC+(stream)|REST&chd=t:259.356,32.297|5.117,3.448|260.499,58.950&chds=a&chxt=x,y&chxl=0:|1x|16x|&chco=4D89F9,FF6600,C6D9FD&chxs=1N**+µs&chma=10,10,10,10&chbh=30,5,20 "gRPC vs REST loopback")

### Remote (One physical machine to another on a 1GB LAN)

```
$ GRPC_REMOTE_ADDR=192.168.🚫.🚫:4443 REST_REMOTE_ADDR=192.168.🚫.🚫:4444 go test -bench Remote -benchtime 10s
goos: darwin
goarch: amd64
pkg: github.com/anzdaddy/grpc_v_rest
BenchmarkGRPCSetInfoRemote-4                    5000       3585969 ns/op
BenchmarkGRPCSetInfoRemote16x-4                50000        271611 ns/op
BenchmarkGRPCSetInfoStreamRemote-4           3000000          4156 ns/op³
BenchmarkGRPCSetInfoStreamRemote16x-4        5000000          3378 ns/op³
BenchmarkRESTSetInfoRemote-4                    5000       3485564 ns/op
BenchmarkRESTSetInfoRemote16x-4                50000        329956 ns/op
PASS
ok      github.com/anzdaddy/grpc_v_rest    109.602s
```

³Hmm, streaming is even faster over the wire!

![](https://chart.googleapis.com/chart?cht=bvg&chs=400x300&chdl=gRPC|gRPC+(stream)|REST&chd=t:3.585969,0.271611|0.004156,0.003378|3.485564,0.329956&chds=a&chxt=x,y,y&chxl=0:|1x|16x|2:||ms|&chco=4D89F9,FF6600,C6D9FD&chma=10,10,10,10&chbh=30,5,20&chxs=1N*s "gRPC vs REST loopback")

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
