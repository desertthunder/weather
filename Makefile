.PHONY: clean build install run test coverage cover record

clean:
	rm -rf ./.bin
	rm *.json

build:
	@go build -o ./.bin/ ./internal/./... ./cmd/cli/./... ./cmd/geocast/./...

install:
	@echo "Installing geocast..."
	@go install ./internal/./... ./cmd/cli/./... ./cmd/geocast/./...
	@echo "Adding geocast to PATH..."
	@asdf reshim
	@echo "Installed geocast! 🎉"
	@echo "Run 'geocast --help' to get started."

run:
	@./.bin/cli

record:
	@vhs assets/demo.tape

test:
	@mkdir -p .cov
	@go test -v ./... -coverprofile=.cov/coverage.out

coverage:
	@go tool cover -html=.cov/coverage.out -o .cov/coverage.html
	@go tool cover -func=.cov/coverage.out | tee .cov/coverage.txt

cover:
	@python coverage.py .cov/coverage.txt
