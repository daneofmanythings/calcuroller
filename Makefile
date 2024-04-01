gen_grpc:
	protoc \
	--go_out=. \
	--go_opt=paths=source_relative \
	--go-grpc_out=. \
	--go-grpc_opt=paths=source_relative \
	./internal/grpc/proto/roller.proto

run:
	go run ./internal/grpc/server/server.go

test:
	go test ./...
