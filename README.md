```
go get google.golang.org/grpc
go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
go get google.golang.org/protobuf/cmd/protoc-gen-go

cd scistream
protoc --go_out=. --go_opt=paths=source_relative        --go-grpc_out=. --go-grpc_opt=paths=source_relative        scistream.proto
```
