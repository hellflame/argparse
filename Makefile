tidy:
	find . -name "*.go" -type f | xargs -n1 go fmt

test:
	go test -v .