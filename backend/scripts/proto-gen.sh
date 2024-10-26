go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

export PATH="$PATH:$(go env GOPATH)/bin"

mkdir -p ./src/proto/pb

protoc --proto_path=./src/proto \
  --go_out=./src/proto/pb \
  --go_opt=paths=source_relative \
  --go-grpc_out=./src/proto/pb \
  --go-grpc_opt=paths=source_relative \
  ./src/proto/bee_management.proto
