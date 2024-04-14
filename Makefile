.PHONY: build
build: clean
	@go build -o dist/mea

.PHONY: clean
clean:
	@rm -fr dist

.PHONY: test
test:
	@go test -v ./...

.PHONY: test-with-update
test-with-update:
	@go test -v -update ./...
