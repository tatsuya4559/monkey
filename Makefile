build:
	go build main.go

.PHONY: test
test:
	go test ./...

.PHONY: fmt
fmt:
	@gofmt -l .

.PHONY: lint
lint:
	staticcheck ./...
