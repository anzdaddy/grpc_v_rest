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
./grpc_v_rest -cert cert.pem -key key.pem
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
./grpc_v_rest -cert cert.pem -key key.pem
```

Machine 2:

```
GRPC_REMOTE_ADDR=<<machine1-ip>>:4443 REST_REMOTE_ADDR=<<machine1-ip>>:4444 go test -bench Remote -benchtime 10s
```

## Results

Loopback:

```
$ go test -v -bench Loopback
goos: darwin
goarch: amd64
pkg: github.com/anzdaddy/grpc_v_rest
BenchmarkGRPCSetInfoLoopback-8          10000        168310 ns/op
BenchmarkRESTSetInfoLoopback-8          10000        160458 ns/op
PASS
ok      github.com/anzdaddy/grpc_v_rest    5.087s
```

Local:

```
$ GRPC_REMOTE_ADDR=localhost:4443 REST_REMOTE_ADDR=localhost:4444 go test -bench Remote
goos: darwin
goarch: amd64
pkg: github.com/anzdaddy/grpc_v_rest
BenchmarkGRPCSetInfoRemote-8          10000        204239 ns/op
BenchmarkRESTSetInfoRemote-8          10000        183406 ns/op
PASS
ok      github.com/anzdaddy/grpc_v_rest    5.341s
```

Remote:

```
$ GRPC_REMOTE_ADDR=192.168.🚫.🚫:4443 REST_REMOTE_ADDR=192.168.🚫.🚫:4444 go test -bench Remote -benchtime 10s
goos: darwin
goarch: amd64
pkg: github.com/anzdaddy/grpc_v_rest
BenchmarkGRPCSetInfoRemote-4           5000       3630114 ns/op
BenchmarkRESTSetInfoRemote-4           5000       3459960 ns/op
PASS
ok      github.com/anzdaddy/grpc_v_rest    36.257s
```

Machine 1 (also used for loopback and local):

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

Machine 2:

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
