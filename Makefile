# Build the current machine
.PHONY: build
build:
	@go build -o bin/commit ./cmd/commit/

# Build for all platforms
.PHONY: all
all: windows linux macos

# Build for MacOS
.PHONY: macos
macos:
	@GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o bin/macos-amd64/commit ./cmd/commit/
	@GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -o bin/macos-arm64/commit ./cmd/commit/

# Build for Linux
.PHONY: linux
linux:
	@GOOS=linux GOARCH=386 go build -ldflags "-s -w" -o bin/linux-386/commit ./cmd/commit/
	@GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o bin/linux-amd64/commit ./cmd/commit/
	@GOOS=linux GOARCH=arm go build -ldflags "-s -w" -o bin/linux-arm/commit ./cmd/commit/
	@GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -o bin/linux-arm64/commit ./cmd/commit/

# Build for Windows
.PHONY: windows
windows:
	@GOOS=windows GOARCH=386 go build -ldflags "-s -w" -o bin/windows-386/commit ./cmd/commit/
	@GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o bin/windows-amd64/commit ./cmd/commit/
	@GOOS=windows GOARCH=arm go build -ldflags "-s -w" -o bin/windows-arm/commit ./cmd/commit/
	@GOOS=windows GOARCH=arm64 go build -ldflags "-s -w" -o bin/windows-arm64/commit ./cmd/commit/

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
