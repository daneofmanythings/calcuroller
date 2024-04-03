gen_grpc:
	protoc \
	--go_out=. \
	--go_opt=paths=source_relative \
	--go-grpc_out=. \
	--go-grpc_opt=paths=source_relative \
	./internal/grpc/proto/roller.proto

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o diceroni .

.PHONY: run
run: build
	./diceroni

.PHONY: docker-build
docker-build:
		docker build . -t diceroni:latest

.PHONY: docker-run
docker-run: docker-build
	docker run -p 8081:8080 diceroni

.PHONY: test
test:
	go test ./...

.PHONY: ping
ping:
	grpcurl -plaintext -format text -d 'dice_string: "d20 + 5", caller_id: "Joe"' localhost:8081 Roller.Roll
