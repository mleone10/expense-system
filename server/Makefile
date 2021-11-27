.PHONY: build run

clean:
	rm -rf bin

build: clean
	for CMD in `ls cmd`; do \
		env GOOS=linux go build -ldflags="-s -w" -o bin/$$CMD ./cmd/$$CMD/...; \
	done

run: build
	./bin/localserver
