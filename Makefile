.PHONY: build
build:
	@go build -o bin/ ./cmd/commit 

.PHONY: clean
clean:
	@if [ -f bin/commit ]; then \
		rm bin/commit ; \
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
