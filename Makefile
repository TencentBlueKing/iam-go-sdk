.PHONY: init dep migrations mock lint lint-dupl test bench build build-linux build-aarch64 clean all serve

VERSION = `head -1 VERSION`

init:
	pip install pre-commit
	pre-commit install
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.27.0

dep:
	go mod tidy

mock:
	go generate ./...

lint:
	golangci-lint run

lint-dupl:
	golangci-lint run --no-config --disable-all --enable=dupl

test:
	go test -gcflags=all=-l $(shell go list ./... | grep -v mock | grep -v docs) -covermode=count -coverprofile .coverage.cov

bench:
	go test -run=nonthingplease -benchmem -bench=. $(shell go list ./... | grep -v /vendor/)
