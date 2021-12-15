.PHONY: build run
.SILENT:
-include ./.env

clean:
	rm -rf bin

build: clean
	for CMD in `ls cmd`; do \
		env GOOS=linux go build -ldflags="-s -w" -o bin/$$CMD ./cmd/$$CMD/...; \
	done

run: build
	# TODO: Find a better way to manage local configuration.
	COGNITO_CLIENT_SECRET=${COGNITO_CLIENT_SECRET} ./bin/localserver
