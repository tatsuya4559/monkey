build:
	go build main.go

.PHONY: test
test:
	go test ./lexer ./ast ./parser ./object ./evaluator

.PHONY: fmt
fmt:
	@gofmt -l .

.PHONY: lint
lint:
	golint ./token ./lexer ./ast ./parser ./object ./evaluator
