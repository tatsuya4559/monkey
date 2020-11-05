build:
	go build main.go

test:
	go test ./lexer ./ast ./parser ./object ./evaluator

fmt:
	gofmt -l .
