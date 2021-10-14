tidy:
	find . -name "*.go" -type f | xargs -n1 go fmt

test:
	go test -v .

.coverage.out: .FORCE
	go test -coverprofile .coverage.out

.FORCE:

cover-report: .coverage.out
	go tool cover -html .coverage.out
