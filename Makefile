.PHONY: run clean test

clean:
	rm -rf ./.bin
	rm *.json

run:
	@go build -o ./.bin/ ./cmd/cli/./... ./internal/./... && ./.bin/cli

record:
	@vhs assets/demo.tape

test:
	@mkdir -p .cov
	@go test -v ./... -coverprofile=.cov/coverage.out

coverage:
	@go tool cover -html=.cov/coverage.out -o .cov/coverage.html
	@go tool cover -func=.cov/coverage.out | tee .cov/coverage.txt

cover:
	@go tool cover -html=.cov/coverage.out
