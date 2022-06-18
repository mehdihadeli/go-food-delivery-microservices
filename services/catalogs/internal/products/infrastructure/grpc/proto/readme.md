- [Compile gRPC Quick start](https://grpc.io/docs/languages/go/quickstart/)
- [Protocol Buffers and GRPC in Go](https://dev.to/karankumarshreds/protocol-buffers-and-grpc-in-go-3eil)

## Compile all proto files:

``` cmd
protoc --go_out=. --go-grpc_out=.  *.proto
```