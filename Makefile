.PHONY: build
build:
	@go build -o bin/commit cmd/commit/main.go
ifeq ($(filter all,$(MAKECMDGOALS)),all)
	@GOOS=darwin GOARCH=amd64 go build -o bin/macos-amd64/commit ./cmd/commit/
	@GOOS=darwin GOARCH=arm64 go build -o bin/macos-arm64/commit ./cmd/commit/
	@GOOS=linux GOARCH=386 go build -o bin/linux-386/commit ./cmd/commit/
	@GOOS=linux GOARCH=amd64 go build -o bin/linux-amd64/commit ./cmd/commit/
	@GOOS=linux GOARCH=arm go build -o bin/linux-arm/commit ./cmd/commit/
	@GOOS=linux GOARCH=arm64 go build -o bin/linux-arm64/commit ./cmd/commit/
	@GOOS=windows GOARCH=386 go build -o bin/windows-386/commit ./cmd/commit/
	@GOOS=windows GOARCH=amd64 go build -o bin/windows-amd64/commit ./cmd/commit/
	@GOOS=windows GOARCH=arm go build -o bin/windows-arm/commit ./cmd/commit/
	@GOOS=windows GOARCH=arm64 go build -o bin/windows-arm64/commit ./cmd/commit/
endif

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
