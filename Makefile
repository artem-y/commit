.PHONY: build
build:
	@go build -o bin/ ./cmd/commit

.PHONY: release
release:
	@GOOS=darwin GOARCH=amd64 go build -o bin/macos/commit -ldflags "-s -w" ./cmd/commit
	@GOOS=linux GOARCH=amd64 go build -o bin/linux-arm64/commit -ldflags "-s -w" ./cmd/commit
	@GOOS=linux GOARCH=386 go build -o bin/commit-linux-x86 -ldflags "-s -w" ./cmd/commit
	@GOOS=windows GOARCH=amd64 go build -o bin/windows-arm64/commit -ldflags "-s -w" ./cmd/commit
	@GOOS=windows GOARCH=386 go build -o bin/windows-x86/commit -ldflags "-s -w" ./cmd/commit

.PHONY: clean
clean:
	@if [ -d bin/ ]; then \
		rm -rf bin ; \
	fi;
	@go clean

.PHONY: run
run:
	@go run ./cmd/commit "$(filter-out $@,$(MAKECMDGOALS))"

.PHONY: test
test:
	@go test -v ./tests/

%:
	@:
