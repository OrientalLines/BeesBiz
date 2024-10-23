go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

export PATH="$PATH:$(go env GOPATH)/bin"

mkdir -p ../proto/pb

protoc --proto_path=../proto \
  --go_out=../proto/pb \
  --go_opt=paths=source_relative \
  --go-grpc_out=../proto/pb \
  --go-grpc_opt=paths=source_relative \
  ../proto/bee_management.proto
