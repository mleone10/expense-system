.PHONY: build run
.SILENT:
-include ./.env

clean:
	rm -rf bin

build: clean
	go vet ./...
	go fmt ./...
	for CMD in `ls cmd`; do \
		env GOOS=linux go build -ldflags="-s -w" -o bin/$$CMD ./cmd/$$CMD/...; \
	done

run: build
	COGNITO_CLIENT_SECRET=${COGNITO_CLIENT_SECRET} ./bin/localserver
