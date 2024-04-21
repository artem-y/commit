# Build the current machine
.PHONY: build
build:
	@go build -o bin/ ./cmd/commit/

# Build for all platforms
.PHONY: all
all: windows linux macos

# Build for MacOS
.PHONY: macos
macos:
	@GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o bin/macos-amd64/ ./cmd/commit/
	@GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -o bin/macos-arm64/ ./cmd/commit/

# Build for Linux
.PHONY: linux
linux:
	@GOOS=linux GOARCH=386 go build -ldflags "-s -w" -o bin/linux-386/ ./cmd/commit/
	@GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o bin/linux-amd64/ ./cmd/commit/
	@GOOS=linux GOARCH=arm go build -ldflags "-s -w" -o bin/linux-arm/ ./cmd/commit/
	@GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -o bin/linux-arm64/ ./cmd/commit/

# Build for Windows
.PHONY: windows
windows:
	@GOOS=windows GOARCH=386 go build -ldflags "-s -w" -o bin/windows-386/ ./cmd/commit/
	@GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o bin/windows-amd64/ ./cmd/commit/
	@GOOS=windows GOARCH=arm go build -ldflags "-s -w" -o bin/windows-arm/ ./cmd/commit/
	@GOOS=windows GOARCH=arm64 go build -ldflags "-s -w" -o bin/windows-arm64/ ./cmd/commit/

# Build and prepare a release archive for all platforms
.PHONY: release
release:
	@make clean
	@make all
	@zip -r bin/bin.zip bin/

# Clean the build artifacts
.PHONY: clean
clean:
	@if [ -d bin/ ]; then \
		rm -rf bin ; \
	fi;
	@go clean

# Build and run the application
.PHONY: run
run:
	@go run ./cmd/commit "$(filter-out $@,$(MAKECMDGOALS))"

# Run the tests
.PHONY: test
test:
	@go test -v ./tests/

%:
	@:
