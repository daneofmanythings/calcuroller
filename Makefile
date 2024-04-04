gen_grpc:
	@ echo Generating grpc with protoc
	@ protoc \
	--go_out=. \
	--go_opt=paths=source_relative \
	--go-grpc_out=. \
	--go-grpc_opt=paths=source_relative \
	./internal/grpc/proto/roller.proto
	@ echo ...done!

.PHONY: build-repl
build-repl:
	@ echo Locally building binary
	@ mkdir .bin/ -p
	@ go build -o repl ./pkg/interpreter/
	@ mv repl ./.bin/
	@ echo ...done!

.PHONY: run-repl
run-repl: build-repl
	@ ./.bin/repl

.PHONY: build-server-local
build-server-local:
	@ echo Locally building binary
	@ mkdir .bin/ -p
	@ go build -o local_server ./internal/grpc/server/server.go
	@ mv local_server .bin/
	@ echo ...done!

.PHONY: run-server-local
run-server-local: build-server-local
	@ ./.bin/local_server

.PHONY: build-server-docker
build-server-docker:
	@ docker build . -t calcuroller:latest

.PHONY: run-server-docker
run-server-docker: build-server-docker
	@ docker run --publish 8080:8080 calcuroller

.PHONY: run-server-docker-host
run-server-docker-host: build-server-docker
	@ docker run --network="host" calcuroller

 .PHONY: build-server-docker-multistage
build-server-docker-multistage:
	@ docker build -f Dockerfile.multistage  ./internal/grpc/server/ -t calcuroller

.PHONY: run-server-docker-multistage
run-server-docker-multistage: build-server-docker-multistage
	@ docker run --publish 8080:8080 calcuroller


.PHONY: run-server-docker-multistage-host
run-server-docker-multistage-host: build-server-docker-multistage
	@ docker run --network="host" calcuroller

.PHONY: test
test:
	@ echo Running tests
	@ go test ./...

.PHONY: ping
ping:
	@ echo Pinging with 'dice_string: "d20 + 5", caller_id: "Joe"'
	@ grpcurl -plaintext -format text -d 'dice_string: "d20 + 5", caller_id: "Joe"' localhost:8080 Roller.Roll

.PHONY: clean
clean:
	@ echo Removing locally compiled files
	@ rm -rf ./.bin/
	@ echo ...done!
